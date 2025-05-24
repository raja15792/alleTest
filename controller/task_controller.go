package controller

import (
	"net/http"
	"time"

	"github.com/alle/tasks/model"
	"github.com/alle/tasks/service"
	"github.com/alle/tasks/common"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TaskController struct {
	manager             service.ITaskManager
	validate            *validator.Validate
}

func NewTaskController(manager service.ITaskManager) *TaskController {
	return &TaskController{
		manager:             manager,
		validate:            validator.New(),
	}
}

// CreateTask creates the task and add a record in postgres
func (r *TaskController) CreateTask(c echo.Context) error {
	var createRequest model.Task

	if err := c.Bind(&createRequest); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	err := r.validate.Struct(createRequest)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	// enrich create request with createdAt and modifiedAt fields
	createRequest.CreatedAt = time.Now().UTC()
	createRequest.ModifiedAt = time.Now().UTC()

	task, code := r.manager.CreateTask(createRequest)
	if code != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, common.PackResponse(code,
			common.StatusText(code), nil))
	}

	return c.JSON(http.StatusOK, common.PackResponse(code, task, nil))
}

// UpdateTask update the task for given task id.
func (r *TaskController) UpdateTask(c echo.Context) error {
	var updateRequest model.Task

	// Bind payload into struct
	if err := c.Bind(&updateRequest); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	// get taskId from path param
	taskId, err := uuid.Parse(c.Param("id"))
	updateRequest.Id = taskId

	// validate update payload
	err = r.validate.Struct(updateRequest)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	
	// update the modified date
	updateRequest.ModifiedAt = time.Now().UTC()

	var code = r.manager.UpdateTask(updateRequest)
	if code != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, common.PackResponse(code,
			common.StatusText(code), nil))
	}

	return c.JSON(http.StatusOK, common.PackResponse(code, updateRequest.Id, nil))
}

// GetAllTasks returns tasks
// If you provide query parameters for filtering, along with page and per-page values, it will return the corresponding results.
func (r *TaskController) GetAllTasks(c echo.Context) error {
	var filter service.TaskFilterParams

	// Bind query params into TaskFilterParams
	if err := c.Bind(&filter); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	var data, count, status = r.manager.GetAllTasks(c, filter)
	var page, perPage int

	// get per page records value
	if filter.PerPage != nil {
		perPage = int(*filter.PerPage)
	}

	// get page number value
	if filter.Page != nil {
		page = int(*filter.Page)
	}

	meta := common.ResponseMeta{
		Total:   count,
		Page:    page,
		PerPage: perPage}
	if status != common.DbError {
		return c.JSON(http.StatusOK, common.PackResponse(http.StatusOK, data, &meta))
	} else {
		return c.JSON(http.StatusInternalServerError, common.PackResponse(status,
			common.StatusText(status), &meta))
	}
}

// DeleteTask update the task for given task id.
func (r *TaskController) DeleteTask(c echo.Context) error {
	// get taskId from path param
	taskId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	
	// delete the record from table
	var code = r.manager.DeleteTask(taskId.String())
	if code != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, common.PackResponse(code,
			common.StatusText(code), nil))
	}

	return c.JSON(http.StatusOK, common.PackResponse(code, taskId, nil))
}