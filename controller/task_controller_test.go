package controller_test

import (
	"bytes"
	"encoding/json"
	"github.com/alle/tasks/controller"
	"github.com/alle/tasks/model"
	"github.com/alle/tasks/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"

)

// MockTaskManager is a mock implementation of the ITaskManager interface
type MockTaskManager struct {
	mock.Mock
}

func (m *MockTaskManager) CreateTask(task model.Task) (*model.TaskId, int) {
	args := m.Called(task)
	return args.Get(0).(*model.TaskId), args.Int(1)
}

func (m *MockTaskManager) UpdateTask(task model.Task) int {
	args := m.Called(task)
	return args.Int(0)
}

func (m *MockTaskManager) DeleteTask(id string) int {
	args := m.Called(id)
	return args.Int(0)
}

func (m *MockTaskManager) GetAllTasks(c echo.Context, filter service.TaskFilterParams) ([]*model.Task, int, int) {
	args := m.Called(c, filter)
	return args.Get(0).([]*model.Task), args.Int(1), args.Int(2)
}


func TestCreateTask_Success(t *testing.T) {
	e := echo.New()
	mockManager := new(MockTaskManager)
	controller := controller.NewTaskController(mockManager)

	task := model.Task{
		Name: "Updated Task",
		Status: "Pending",
	}
	body, _ := json.Marshal(task)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedID := &model.TaskId{Id: uuid.New()}
	mockManager.On("CreateTask", mock.AnythingOfType("model.Task")).Return(expectedID, http.StatusOK)

	err := controller.CreateTask(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockManager.AssertExpectations(t)
}


func TestUpdateTask_Success(t *testing.T) {
	e := echo.New()
	mockManager := new(MockTaskManager)
	controller := controller.NewTaskController(mockManager)

	taskId := uuid.New()
	task := model.Task{
		Name: "Updated Task",
		Status: "Pending",
	}
	body, _ := json.Marshal(task)

	req := httptest.NewRequest(http.MethodPut, "/tasks/"+taskId.String(), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(taskId.String())

	mockManager.On("UpdateTask", mock.AnythingOfType("model.Task")).Return(http.StatusOK)

	err := controller.UpdateTask(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAllTasks_Success(t *testing.T) {
	e := echo.New()
	mockManager := new(MockTaskManager)
	controller := controller.NewTaskController(mockManager)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTasks := []*model.Task{
		{Name: "Task1", Status: "Pending"},
		{Name: "Task2", Status: "Pending"},
	}
	mockManager.On("GetAllTasks", c, mock.AnythingOfType("service.TaskFilterParams")).
		Return(mockTasks, len(mockTasks), http.StatusOK)

	err := controller.GetAllTasks(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockManager.AssertExpectations(t)
}


func TestDeleteTask_Success(t *testing.T) {
	e := echo.New()
	mockManager := new(MockTaskManager)
	controller := controller.NewTaskController(mockManager)

	taskId := uuid.New()

	req := httptest.NewRequest(http.MethodDelete, "/tasks/"+taskId.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(taskId.String())

	mockManager.On("DeleteTask", taskId.String()).Return(http.StatusOK)

	err := controller.DeleteTask(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
