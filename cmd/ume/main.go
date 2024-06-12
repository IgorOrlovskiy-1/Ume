package main

import (
    "Ume/internal/config"
    "Ume/components/hello_templ"
	"fmt"

    "github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Handler
func nullpage(c echo.Context) error {
    component := hello_templ.hello("World")
    return component
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
	e.GET("/", templ.Handler(nullpage))

	// Start server
    e.Logger.Fatal(e.Start(":" + cfg.Port))

	//TODO: tests
}
