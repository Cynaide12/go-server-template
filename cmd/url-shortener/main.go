package main

import (
	"log/slog"
	"net/http"
	"os"
	"url_shortener/internal/config"
	httphandlers "url_shortener/internal/http-server/handlers"
	"url_shortener/internal/http-server/logger"
	slogpretty "url_shortener/internal/lib/logger/handlers"
	"url_shortener/internal/lib/logger/sl"
	"url_shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	log.Debug("debug messages are enabled")

	log.Debug("start init storage")
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	log.Debug("storage init succesful")

	initRouter(cfg, log, storage)
}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = setupPrettySlog()
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log

}

func initRouter(cfg *config.Config, log *slog.Logger, UrlHandler httphandlers.UrlHandler) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(logger.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/{alias}", httphandlers.RedirectURLHandler(log, UrlHandler))

	r.Post("/new-url", httphandlers.SaveURLHandler(log, UrlHandler))
	r.Post("/delete-url", httphandlers.DeleteURLHandler(log, UrlHandler))

	r.Get("/get-url", httphandlers.GetURLHandler(log, UrlHandler))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("starting server", slog.String("address", srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))

		os.Exit(1)
	}

	log.Error("server stopped")

}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
