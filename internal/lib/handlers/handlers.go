package handlers

import (
	"Ume/components"
	"Ume/internal/config"
	"Ume/internal/lib/logger/sl"
	"fmt"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"net/http"
)

var cfg = config.MustLoad()
var log = sl.SetupLogger(cfg.Env)

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
