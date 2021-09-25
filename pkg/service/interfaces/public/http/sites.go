package http

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"

	app "github.com/waffleboot/playstation_buy/pkg/service/application"
)

type Update struct {
	Host  string `json:"host"`
	Count int    `json:"count"`
}

type Endpoint struct {
	service *app.Service
}

func NewEndpoint(s *app.Service) *Endpoint {
	return &Endpoint{service: s}
}

func (e *Endpoint) sites(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	m, err := e.service.ProcessQuery(search)
	if errors.Is(err, context.DeadlineExceeded) {
		w.WriteHeader(http.StatusRequestTimeout)
		render.JSON(w, r, m)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, m)
		return
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, m)
}

func (e *Endpoint) update(w http.ResponseWriter, r *http.Request) {
	var req Update
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	e.service.Update(req.Host, req.Count)
	w.WriteHeader(http.StatusOK)
}
