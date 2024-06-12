package main

import (
	"Ume/internal/config"
	"Ume/internal/lib/logger/sl"
	"Ume/internal/storage/postgresql"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	db, err := postgresql.New(
		"user=" + cfg.User + " password=" + cfg.Password + " dbname=" + cfg.DBName,
	)
	if err != nil {
		log.Error("Failed to connect to database", sl.Err(err))
		os.Exit(1)
	}
	log.Info("successfully connected to database")

	//test
	id, err := db.AddUser("tester", "testerov", "12345")
	if err != nil {
		log.Error("Failed to add user", sl.Err(err))
		os.Exit(1)
	}
	log.Info("successfully add user to table users")

	_ = id
	//TODO: router

	//TODO: middlewars

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
