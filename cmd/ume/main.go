package main

import (
	"Ume/components"
	"Ume/internal/config"
	"Ume/internal/lib/logger/sl"
	"Ume/internal/storage/postgresql"
	"net/http"

	"fmt"
	"log/slog"
	"os"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var cfg = config.MustLoad()
var log = setupLogger(cfg.Env)

func CheckUserCreds(login string, password string) (error, bool) {
	// TODO
	if login == "123" && password == "123" {
		return nil, true
	}
	return nil, false
}

func Login(c echo.Context) error {
	login := c.FormValue("login")
	password := c.FormValue("password")
	err, status := CheckUserCreds(login, password)
	if err != nil {
		log.Info("Failed to login", sl.Err(err))
	}
	fmt.Printf("Status:%t Login:%s  Password:%s\n", status, login, password)
	if !status {
        return Render(c, http.StatusOK, components.BadLogin())
	}
	return Render(c, http.StatusOK, components.GoodLogin())
}

func HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, components.Home())
}

func ChatHandler(c echo.Context) error {
	return Render(c, http.StatusOK, components.Chat())
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func main() {
	db, err := postgresql.New(
		"user=" + cfg.User + "password=" + cfg.Password + "dbname=" + cfg.DBName + "ssl=" + cfg.SSLMode,
	)
	if err != nil {
		log.Error("Failed to connect to database", sl.Err(err))
		//os.Exit(1)
	}

	_ = db

	// Server instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/assets", "components/assets")
	e.GET("/", HomeHandler)
	e.GET("/chat", ChatHandler)
	e.POST("/login", Login)
	// Start server
	s := http.Server{
		Addr:        ":" + cfg.Port,
		Handler:     e,
		ReadTimeout: cfg.Timeout,
	}
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	}
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
