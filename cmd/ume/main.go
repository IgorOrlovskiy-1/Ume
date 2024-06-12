package main

import (
	"Ume/internal/config"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	logger := middleware.Logger()

	//TODO: bd

	// Server instance
	e := echo.New()

	// Middleware
	e.Use(logger)
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":1024"))

	//TODO: tests
}
