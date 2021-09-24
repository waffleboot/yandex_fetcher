package http

import (
	"github.com/go-chi/chi"
)

func AddRoutes(router *chi.Mux) {
	router.Get("/sites", sites)
}
