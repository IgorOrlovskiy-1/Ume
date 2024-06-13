package main

import (
	"Ume/internal/config"
	"Ume/internal/lib/handlers"
	"Ume/internal/lib/logger/sl"
	"Ume/internal/storage/postgresql"
	"net/http"

	_ "os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var cfg = config.MustLoad()
var log = sl.SetupLogger(cfg.Env)

func main() {
	db, err := postgresql.New(
		"user=" + cfg.User + " password=" + cfg.Password + " dbname=" + cfg.DBName,
	)
    _ = db
	if err != nil {
		log.Error("Failed to connect to database", sl.Err(err))
		//os.Exit(1)
	}
	log.Info("successfully connected to database")

	//test
	//id, err := db.AddUser("tester", "testerov", "12345")
	if err != nil {
		log.Error("Failed to add user", sl.Err(err))
		//os.Exit(1)
	}
	log.Info("successfully add user to table users")

	//_ = id

	// Server instance
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//TODO: router
	e.Static("/assets", "components/assets")
	e.GET("/", handlers.HomeHandler)
	e.GET("/chat", handlers.ChatHandler)
	e.POST("/login", handlers.Login)
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
