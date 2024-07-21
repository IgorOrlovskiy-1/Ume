package userCreate

import (
	"errors"
	"net/http"
	"log/slog"
	//"time"
	resp "Ume/internal/lib/api/response"
	"Ume/internal/lib/logger/sl"
	"Ume/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AddNewUser interface {
	AddUser(firstName, lastName, password, email, username string) (error)
}

type Request struct {
	FirstName string `json:"first_name" validate:"required"`
	SecondName string `json:"second_name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	// DateBirthday time.Time `json:"date_birthday"`
}

type Response struct {
	resp.Response
}

func NewUser(log *slog.Logger, newUser AddNewUser) http.HandlerFunc	{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.create.new.user"
		
		log = log.With (
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

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
		
		password, err := hashPassword(req.Password)
		if err != nil {
			log.Error("Error with hashing password, try again", sl.Err(err))

			render.JSON(w, r, resp.Error("Error with hashing password, try again"))
		}

		err = newUser.AddUser(req.FirstName, req.SecondName, password, req.Email, req.Username)
		if errors.Is(err, storage.ErrUserWithUsernameExists) {
			log.Info("Username already exists", slog.String("username", req.Username))

			render.JSON(w, r, resp.Error("Username already exists"))

			return
		}
		if err != nil {
			log.Error("Error with adding new user", sl.Err(err))

			render.JSON(w, r, resp.Error("Error with adding new user"))

			return
		}

		log.Info("User added", slog.Any("username", req.Username))

		render.JSON(w, r, Response{
			Response: resp.OK()}) 
	}
}

func hashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func checkPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}