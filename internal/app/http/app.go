package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/exp/slog"
	"url-validator/internal/http-server/handler/url"
	mwLogger "url-validator/internal/http-server/middleware/logger"
	sUrl "url-validator/internal/service/url"
)

type App struct {
	Log    *slog.Logger
	Router *chi.Mux
}

func New(log *slog.Logger) *App {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	urlService := sUrl.New(log)
	urlApi := url.New(log, urlService)

	router.Route("/api", func(r chi.Router) {
		r.Post("/urls/validate", urlApi.Validate())
	})

	router.Mount("/swagger", httpSwagger.WrapHandler)

	return &App{
		Log:    log,
		Router: router,
	}
}
