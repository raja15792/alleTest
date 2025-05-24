package main

import (
	"fmt"
	"log"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/caarlos0/env/v6"
	"github.com/alle/tasks/common"
	"github.com/alle/tasks/service"
	"github.com/alle/tasks/controller"
	"github.com/alle/tasks/db"
)

func main() {
	conf := &common.Config{}
	err := env.Parse(conf)
	if err != nil {
		log.Panic(err)
	}

	log.Println("config loaded")
	// Echo instance -> Using echo as a http routing 
	e := echo.New()

	log.Println("loading db")
	log.Println(conf.DbURI)
	psql, err := db.NewPgPool(conf.DbURI)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	if err = psql.Ping(context.Background()); err != nil {
		log.Println(err.Error())
		panic(err)
	}
	log.Println("db initialized")

	taskService := service.NewTaskManager(psql)
	taskController := controller.NewTaskController(taskService)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 2,
	}))

	// Heathcheck endpoint
	hc := e.Group("/healthcheck")
	hc.GET("/alive", ok)

	// Task manager endpoint
	task := e.Group("/v1")
	task.POST("/task", taskController.CreateTask)
	task.PATCH("/task/:id", taskController.UpdateTask)
	task.GET("/tasks", taskController.GetAllTasks)
	task.DELETE("/task/:id", taskController.DeleteTask)

	log.Println("starting server...")

	e.Logger.Panic(e.Start(fmt.Sprintf(":%s", "8080")))
}

func ok(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}