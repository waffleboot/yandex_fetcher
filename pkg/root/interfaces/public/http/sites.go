package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/render"

	app "github.com/waffleboot/playstation_buy/pkg/root/application"
)

type Endpoint struct {
	service *app.Service
}

func NewEndpoint(s *app.Service) *Endpoint {
	return &Endpoint{service: s}
}

func (e *Endpoint) sites(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	log.Printf("sites/%s", search)
	m, err := e.service.ProcessQuery(search)
	if errors.Is(err, context.DeadlineExceeded) {
		w.WriteHeader(http.StatusRequestTimeout)
		render.JSON(w, r, m)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, m)
}
