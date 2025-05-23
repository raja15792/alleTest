package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance -> Using echo as a http routing 
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 2,
	}))

	hc := e.Group("/healthcheck")
	hc.GET("/alive", ok)

	log.Println("starting server...")

	e.Logger.Panic(e.Start(fmt.Sprintf(":%s", "8080")))
}

func ok(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}