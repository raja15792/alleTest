package model

import (
	"time"
	"fmt"
	"strings"
	"github.com/google/uuid"
)

type Task struct {
	Id uuid.UUID `json:"id" db:"id"`
	Name string `json:"name" validate:"required" db:"task_name"`
	Status string `json:"status" validate:"required" db:"task_status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
}

type TaskId struct {
	Id uuid.UUID
}

func (t *Task) ToUpdateSQL() (string, []interface{}) {
	var (
		rulesSetCondition []string
		count             = 1
		params            []interface{}
	)

	if t.Name != "" {
		params = append(params, t.Name)
		rulesSetCondition = append(rulesSetCondition, fmt.Sprintf("task_name = $%d", count))
		count++
	}

	if t.Status != "" {
		params = append(params, t.Status)
		rulesSetCondition = append(rulesSetCondition, fmt.Sprintf("task_status = $%d", count))
		count++
	}

	if !t.ModifiedAt.IsZero() {
		params = append(params, t.ModifiedAt)
		rulesSetCondition = append(rulesSetCondition, fmt.Sprintf("modified_at = $%d", count))
		count++
	}

	params = append(params, t.Id)

	return fmt.Sprintf("UPDATE task SET %s WHERE id = $%d", strings.Join(rulesSetCondition, ", "), count), params
}
