package http

import (
	"github.com/go-chi/chi"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	app "github.com/waffleboot/yandex_fetcher/pkg/checker/application"
)

func (e *Endpoint) AddRoutes(router *chi.Mux) {
	router.Post("/check", e.check)
}

type Endpoint struct {
	service *app.Service
}

func NewEndpoint(s *app.Service) *Endpoint {
	return &Endpoint{service: s}
}

type Request struct {
	Host string `json:"host"`
	Url  string `json:"url"`
}

func (e *Endpoint) check(w http.ResponseWriter, r *http.Request) {
	var req Request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "unable to read request: %v", err)
		return
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "unable to close request: %v", err)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "unable to parse request: %v", err)
		return
	}
	count, err := e.service.Benchmark(req.Host, req.Url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to benchmark request: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", count)
}
