package homePage

import (
    "html/template"
    "net/http"
	"log/slog"
	"github.com/go-chi/render"
	"github.com/go-chi/chi/v5/middleware"
	resp "Ume/internal/lib/api/response"
	"Ume/internal/lib/logger/sl"

)

func GetHomePage(log *slog.Logger) http.HandlerFunc	{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.homePage.GetHomePage"

		log = log.With (
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)


		w.Header().Set("Content-Type", "text/html")

		tmpl := template.New("homePage")
		tmpl, err := tmpl.ParseFiles("/home/iorlovskiy/Ume/templates/home.html")
		if err != nil {
			log.Error("Failed to parse home html", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to parse home html"))

			return
		}
	
		err = tmpl.ExecuteTemplate(w, "home.html", nil)
		if err != nil {
			log.Error("Error executing template", sl.Err(err))

			render.JSON(w, r, resp.Error("Error executing template"))

			return
		}
	}
}
