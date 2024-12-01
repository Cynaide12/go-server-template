package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url_shortener/internal/lib/api"
	"url_shortener/internal/lib/logger/sl"
	"url_shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func GetURLHandler(log *slog.Logger, UrlHandler UrlHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const fp = "handlers.httphandlers.GetURLHandler"

		log.With(
			slog.String("fp", fp),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		
		alias := r.URL.Query().Get("alias")

		if alias == "" {
			log.Error("alias parameter required")

			render.JSON(w, r, resp.Error("alias parameter required"))

			return
		}

		log.Info("url query parsed", slog.Any("request", alias))

		url, err := UrlHandler.GetUrl(alias)

		if errors.Is(storage.ErrURLNotFound, err) {
			log.Info("url not found", sl.Err(err))

			render.JSON(w, r, resp.Error("URL not found"))

			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get url"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			URL:      url,
		})

	}
}
