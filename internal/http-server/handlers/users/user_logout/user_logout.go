package userLogout

import (
	"net/http"
	"log/slog"
	resp "Ume/internal/lib/api/response"
	"Ume/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/boj/redistore"

)

type Request struct {
	Username string `json:"username" validate:"required"`
	SessionId string `json:"session_id" validate:"required"`
}

type Response struct {
	resp.Response
}

type IsFoundUser interface {
	FindUserPassword(username string) (string, error)
	GetUserIdByUsername(username string) (int64, error)
}

func LogoutUser(log *slog.Logger, store *redistore.RediStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.logout.logout.user"
		
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

		session, err := store.Get(r, req.SessionId)
        if err != nil {
            log.Error("Error with get session from redis", sl.Err(err))

			render.JSON(w, r, resp.Error("InternalServerError"))

			return
        }

        session.Options.MaxAge = -1
   		err = session.Save(r, w)
    	if err != nil {
            log.Error("Error with deleting session from redis", sl.Err(err))

			render.JSON(w, r, resp.Error("Error with logout"))

        	return
    	}

		render.JSON(w, r, Response{ Response: resp.OK()}) 
	}
}
