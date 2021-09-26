package http

import (
	"fmt"

	"github.com/go-chi/chi"

	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"

	app "github.com/waffleboot/yandex_fetcher/pkg/service/application"
)

type Endpoint struct {
	service *app.Service
}

func NewEndpoint(s *app.Service) *Endpoint {
	return &Endpoint{service: s}
}

func (e *Endpoint) AddRoutes(router *chi.Mux) {
	router.Get("/sites", e.sites)
	router.Get("/yandex", e.yandex)
	router.Post("/update", e.update)
}

func (e *Endpoint) sites(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("search query param is absent"))
		return
	}
	m, err := e.service.ProcessQuery(search)
	if errors.Is(err, context.DeadlineExceeded) {
		w.WriteHeader(http.StatusRequestTimeout)
		render.JSON(w, r, m)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, m)
}

func (e *Endpoint) yandex(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("search query param is absent"))
		return
	}
	m, err := e.service.YandexItems(search)
	if errors.Is(err, context.DeadlineExceeded) {
		w.WriteHeader(http.StatusRequestTimeout)
		render.JSON(w, r, m)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, m)
}

type CacheUpdate struct {
	Host  string `json:"host"`
	Count int    `json:"count"`
}

func (e *Endpoint) update(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "unable to read request: %v", err)
		return
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "unable to close body request: %v", err)
		return
	}
	var req CacheUpdate
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "unable to parse request: %v", err)
		return
	}
	e.service.CacheUpdate(req.Host, req.Count)
	w.WriteHeader(http.StatusOK)
}
