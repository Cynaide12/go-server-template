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
	"github.com/go-playground/validator/v10"
)


func DeleteURLHandler(log *slog.Logger, UrlHandler UrlHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const fp = "handlers.httphandlers.DeleteURLHandler"

		log.With(
			slog.String("fp", fp),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req DeleteRequest

		if err := render.Decode(r, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return

		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))


			render.JSON(w, r, resp.ValidationError(validateErr))
		}

		err := UrlHandler.DeleteURL(req.Alias)

		if errors.Is(storage.ErrURLNotFound, err) {
			log.Info("url not found", sl.Err(err))

			render.JSON(w, r, resp.Error("URL not found"))

			return
		}

		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to delete url"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}