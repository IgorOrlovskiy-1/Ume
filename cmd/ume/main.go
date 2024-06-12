package main

import (
    "Ume/internal/storage/postgresql"
    "Ume/components/hello_templ"
    "log/slog"
    "os"
	"fmt"

    "github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Handler
func nullpage(c echo.Context) error {
    component := hello_templ.hello("World")
    return component
}

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	db, err := postgresql.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to connect to database", sl.Err(err))
		os.Exit(1)
	}

	_ = db

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

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}

	return log
}
