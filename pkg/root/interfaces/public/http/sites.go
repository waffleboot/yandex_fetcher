package http

import (
	"net/http"

	"github.com/go-chi/render"
)

func (e *Endpoint) sites(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	m, err := e.service.Query(search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, m)
}
