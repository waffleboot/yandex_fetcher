package http

import (
	"github.com/go-chi/chi"
)

func (e *Endpoint) AddRoutes(router *chi.Mux) {
	router.Get("/sites", e.sites)
	router.Post("/update", e.update)
}
