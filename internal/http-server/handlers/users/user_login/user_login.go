package userLogin

import (
	"errors"
	"net/http"
	"log/slog"
	resp "Ume/internal/lib/api/response"
	"Ume/internal/lib/logger/sl"
	"Ume/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"github.com/boj/redistore"

)

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	SessionName string `json:"session_name"` 
}

type IsFoundUser interface {
	FindUserPassword(username string) (string, error)
	GetUserIdByUsername(username string) (int64, error)
}

func LoginUser(log *slog.Logger, isUserExists IsFoundUser, store *redistore.RediStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.login.login.user"
		
		log = log.With (
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to decode request body"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		password, err := isUserExists.FindUserPassword(req.Username)
		if errors.Is(err, storage.ErrUserWithUsernameExists) {
			log.Info("User not exist", slog.String("username", req.Username))

			render.JSON(w, r, resp.Error("User not exist, try again"))

			return
		}
		if err != nil {
			log.Error("Error with finding user", sl.Err(err))

			render.JSON(w, r, resp.Error("Error with finding user, try again"))

			return
		}

		if(!checkPassword(password, req.Password)) {
			log.Info("Incorrect password", slog.String("password", req.Password))

			render.JSON(w, r, resp.Error("Incorrect password, try again"))

			return
		}

		log.Info("User founded", slog.Any("username", req.Username))

		sessionName := uuid.New().String()		
		session, err := store.Get(r, sessionName)
        if err != nil {
            log.Error("Error with add session to redis", sl.Err(err))

			render.JSON(w, r, resp.Error("InternalServerError"))

			return
        }

		id, err := isUserExists.GetUserIdByUsername(req.Username)
		if err != nil {
            log.Error("Error with getting user id with username", sl.Err(err))

			render.JSON(w, r, resp.Error("Error with getting user id with username"))

			return
        }

		session.Values["authenticated"] = true
        session.Values["userId"] = id
        err = session.Save(r, w)
        if err != nil {
            log.Error("Error with add values to session", sl.Err(err))

			render.JSON(w, r, resp.Error("InternalServerError"))

			return
		}

		render.JSON(w, r, Response{ Response: resp.OK(), SessionName: sessionName }) 
	}
}

func checkPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}