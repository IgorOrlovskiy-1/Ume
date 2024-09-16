package main

import (
	"Ume/internal/config"
	"Ume/internal/lib/logger/sl"
	"Ume/internal/storage/postgresql"
	// "Ume/internal/storage/redis"
	"Ume/internal/http-server/handlers/users/user_create"
	"Ume/internal/http-server/handlers/users/user_login"
	"Ume/internal/http-server/handlers/users/user_logout"
	mwLogger "Ume/internal/middlewares/logger"
	"log/slog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
    "github.com/boj/redistore"
	"os"
	//"context"
	"net/http"
	//"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	db, err := postgresql.NewPool(
		"user=" + cfg.User + " password=" + cfg.Password + " dbname=" + cfg.DBName,
	)
	if err != nil {
		log.Error("Failed to connect to database", sl.Err(err))
		os.Exit(1)
	}

	log.Info("Successfully connected to database")

	// ctx := context.Background()
	// redisClient, err := redis.NewRedisClient(ctx)
	// if err != nil {
	// 	log.Error("Failed to connect to redis", sl.Err(err))
	// 	os.Exit(1)
	// }

	store, err := redistore.NewRediStore(10, "tcp", ":6379", "", []byte(cfg.RedisStoreSecret))
    if err != nil {
        log.Error("Failed to add redis store", sl.Err(err))
		os.Exit(1)
    }
    defer store.Close()

	//test
	// date := time.Now()
	// err = db.AddUser("igor1", "orlovskiy1", "i1@test.io", "12345", "igorO1", date)
	// if err != nil {
	// 	log.Error("Failed to add user", sl.Err(err))
	// 	os.Exit(1)
	// }
	// log.Info("Successfully add user to table users")


	// err = db.AddFriend("igor", "tester")
	// if err != nil {
	// 	log.Error("Failed to add friend", sl.Err(err))
	// 	os.Exit(1)
	// }
	// log.Info("Successfully add friend to table friends")

	// err = db.AddMessage("igorO", "igorO1", "залупапенисхердавалка хуй блядина !ЗЩ)")
	// if err != nil {
	// 	log.Error("Failed to add message", sl.Err(err))
	// 	os.Exit(1)
	// }
	// log.Info("Successfully add message to table messages")
	

	//TODO: router

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.NewLogger(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//router.Get("/", homePage.GetHomePage(log))
	router.Post("/register", userCreate.NewUser(log, db))
	router.Post("/login", userLogin.LoginUser(log, db, store))
	router.Post("/logout", userLogout.LogoutUser(log, store))
	//router.Post("/new_chat", chatCreate.NewChat())

	log.Info("starting server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:	cfg.Address,
		Handler: router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.IdleTimeout,
	} 
	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("Server stopped")

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
