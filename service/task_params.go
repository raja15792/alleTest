package service

import (
	"fmt"
	"strings"
)

type TaskFilterParams struct {
	Status         *string `query:"task_status"`
	Id             *string `query:"id"`
	Page           *int64  `query:"page"`
	PerPage        *int64  `query:"per_page"`
	SortBy         *string `query:"sort"`
	SortOrder      *string `query:"order"`
}

func (t *TaskFilterParams) GetSorts() string {
	var (
		validOrderKeys = map[string]bool{
			"asc":  true,
			"desc": true,
		}
		validSortKeys = map[string]bool{
			"id": true,
			"created_at": true,
		}
	)
	if t.SortBy == nil {
		return ""
	}

	if _, ok := validSortKeys[*t.SortBy]; !ok {
		return ""
	}

	if t.SortOrder != nil {
		if _, ok := validOrderKeys[*t.SortOrder]; ok {
			return fmt.Sprintf("%s %s", *t.SortBy, *t.SortOrder)
		}
	}

	return *t.SortBy
}

func (t *TaskFilterParams) GetLimit() string {
	var (
		page    int64
		perPage int64
	)

	if t.Page != nil {
		page = *t.Page
	}

	if t.PerPage != nil {
		perPage = *t.PerPage
	}

	if perPage != 0 && page == 0 {
		return fmt.Sprintf("limit %d", perPage)
	}

	if perPage != 0 && page != 0 {
		return fmt.Sprintf("limit %d offset %d", *t.PerPage, (*t.Page-1)**t.PerPage)
	}

	return ""
}

func (t *TaskFilterParams) ToSQLClause() (string, []interface{}) {
	var (
		clauses []string
		params  []interface{}
		count   = 1
	)

	if t.Id != nil {
		clauses = append(clauses, fmt.Sprintf("t.id = $%d", count))
		params = append(params, *t.Id)
		return strings.Join(clauses, " AND "), params
	}

	if t.Status != nil {
		clauses = append(clauses, fmt.Sprintf("t.task_status = $%d", count))
		params = append(params, *t.Status)
		return strings.Join(clauses, " AND "), params
	}

	return strings.Join(clauses, " AND "), params
}
