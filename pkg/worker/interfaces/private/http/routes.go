package http

import "github.com/go-chi/chi"

func (e *Endpoint) AddRoutes(router *chi.Mux) {
	router.Post("/check", e.check)
}
