package app

import (
	"golang.org/x/exp/slog"
	"url-validator/internal/app/http"
)

type App struct {
	HttpApp *http.App
}

func New(log *slog.Logger) *App {
	return &App{
		HttpApp: http.New(log),
	}
}
