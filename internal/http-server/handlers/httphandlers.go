package httphandlers

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url_shortener/internal/lib/api"
	"url_shortener/internal/lib/logger/sl"
	"url_shortener/internal/lib/random"
	"url_shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type UrlSaver interface {
	SaveUrl(urlToSave string, alias string) error
	GetUrl(alias string) (string, error)
	DeleteURL(alias string) error
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: в конфиг закинуть
const aliasLength = 8

func NewURL(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fp = "handlers.httphandlers.New"

		log.With(
			slog.String("fp", fp),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.Decode(r, req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		if req.Alias == "" {
			req.Alias = random.RandomString(aliasLength)
		}

		err = urlSaver.SaveUrl(req.URL, req.Alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("URL already exists", sl.Err(err))

			render.JSON(w, r, resp.Error("URL already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save url"))

			return
		}

		log.Info("url added", slog.String("alias", req.Alias))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    req.Alias,
		})
	}

}

