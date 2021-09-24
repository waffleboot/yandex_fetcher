package http

import (
	"github.com/go-chi/chi"
	app "github.com/waffleboot/playstation_buy/pkg/root/application"
)

type Endpoint struct {
	service app.Service
}

func NewEndpoint(s app.Service) Endpoint {
	return Endpoint{service: s}
}

func (e Endpoint) AddRoutes(router *chi.Mux) {
	router.Get("/sites", e.sites)
}
