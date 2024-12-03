package handlers

import resp "url_shortener/internal/lib/api"

//go:generate go run github.com/vektra/mockery/v2@v2.49.1 --name=UrlHandler
type UrlHandler interface {
	SaveUrl(urlToSave string, alias string) error
	GetUrl(alias string) (string, error)
	DeleteURL(alias string) error
}

type SaveRequest struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type GetRequest struct {
	Alias string `json:"alias" validate:"required"`
}

type DeleteRequest struct {
	Alias string `json:"alias" validate:"required"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
	URL   string `json:"url,omitempty"`
}

// TODO: в конфиг закинуть
const aliasLength = 8
