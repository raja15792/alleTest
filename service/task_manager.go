package service

import (
	"fmt"
	"net/http"
	"context"
	"log"
	"github.com/alle/tasks/model"
	"github.com/alle/tasks/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/google/uuid"

)

type TaskManager struct {
	db  *pgxpool.Pool
	ctx context.Context
}

type ITaskManager interface {
	CreateTask(params model.Task) (*model.TaskId, int)
	UpdateTask(params model.Task) int
	GetAllTasks(c echo.Context, filter TaskFilterParams) ([]*model.Task, int, int)
	DeleteTask(taskID string) int
}

func NewTaskManager(db *pgxpool.Pool) *TaskManager {
	return &TaskManager{
		db:  db,
		ctx: context.Background(),
	}
}

// CreateTask creates a task and push it in db
func (t *TaskManager) CreateTask(createReq model.Task) (*model.TaskId, int) {

	trans, err := t.db.BeginTx(t.ctx, pgx.TxOptions{})

	if err != nil {
		log.Println(err)
		return nil, common.DbError
	}

	query := `INSERT INTO task
		(task_name, task_status, created_at, modified_at)
	VALUES
		($1, $2, $3, $4)
	RETURNING id;`

	row := trans.QueryRow(t.ctx, query, createReq.Name, createReq.Status, createReq.CreatedAt, createReq.ModifiedAt)

	var id uuid.UUID
	if err = row.Scan(&id); err != nil {
		trans.Rollback(t.ctx)
		log.Println(err)
		return nil, common.DbError
	}

	err = trans.Commit(t.ctx)
	if err != nil {
		return nil, common.DbError
	}

	return &model.TaskId{Id: id}, http.StatusOK
}

// UpdateTask Update Task
func (t *TaskManager) UpdateTask(updateReq model.Task) int {

	trans, err := t.db.BeginTx(t.ctx, pgx.TxOptions{})

	if err != nil {
		log.Println(err)
		return common.DbError
	}

	// TODO Optimization:  we can write a logic here to check if we are updating the same version of data

	updateSql, params := updateReq.ToUpdateSQL()
	_, err = trans.Exec(t.ctx, updateSql, params...)
	if err != nil {
		trans.Rollback(t.ctx)
		log.Println(err)
		return common.DbError
	}
	err = trans.Commit(t.ctx)
	if err != nil {
		return common.DbError
	}

	return http.StatusOK
}


// GetAllTasks return all tasks. It accept taskFilterParams which 
func (t *TaskManager) GetAllTasks(c echo.Context, filter TaskFilterParams) ([]*model.Task, int, int) {
	var baseQuery string
	baseQuery = `SELECT id, task_name, task_status, created_at, modified_at from task t `
	countQuery := `SELECT count(*) from task t`
	whereClause, values := filter.ToSQLClause()
	sortParams := filter.GetSorts()
	limitClause := filter.GetLimit()

	// check if whereClause is empty -> means there is no where clause condition in query
	if whereClause != "" {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, whereClause)
		countQuery = fmt.Sprintf("%s WHERE %s", countQuery, whereClause)
	}

	// Append order by here
	if sortParams != "" {
		baseQuery = fmt.Sprintf("%s ORDER BY %s", baseQuery, sortParams)
	}

	// check if limitClause is empty -> means if there is no limit condition in query (it will 
	// return all tasks in same 1st page)
	if limitClause != "" {
		baseQuery = fmt.Sprintf("%s %s", baseQuery, limitClause)
	}

	var paramsResponse []*model.Task
	err := pgxscan.Select(context.Background(), t.db, &paramsResponse, baseQuery, values...)

	if err != nil {
		log.Println(err)
		return nil, 0, common.DbError
	}

	var count []int
	err = pgxscan.Select(context.Background(), t.db, &count, countQuery, values...)

	if err != nil {
		log.Println(err)
		return nil, 0, common.DbError
	}
	return paramsResponse, count[0], http.StatusOK
}

// Delete Task given task id
func (t *TaskManager) DeleteTask(taskID string) int {
	// get uuid taskId
	id, err := uuid.Parse(taskID)
	if err != nil {
		log.Println(err)
		return common.DbError
	}

	// Begin transaction block
	trans, err := t.db.BeginTx(t.ctx, pgx.TxOptions{})
	if err != nil {
		log.Println(err)
		return common.DbError
	}

	// initialize delete query var
	query := `DELETE FROM task WHERE id = $1;`

	_, err = trans.Exec(t.ctx, query, id)
	if err != nil {
		trans.Rollback(t.ctx)
		log.Println(err)
		return common.DbError
	}

	err = trans.Commit(t.ctx)
	if err != nil {
		log.Println(err)
		return common.DbError
	}

	return http.StatusOK
}
