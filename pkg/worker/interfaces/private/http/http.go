package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	app "github.com/waffleboot/playstation_buy/pkg/worker/application"
)

type channelItem struct {
	req  Request
	done chan domain.StatsItem
	errc chan error
}

type Endpoint struct {
	channel chan channelItem
}

func NewEndpoint(s *app.Service) *Endpoint {
	channel := make(chan channelItem, 1)
	go func() {
		go func() {
			for e := range channel {
				data, err := s.Benchmark(domain.YandexItem{
					Host: e.req.Host,
					Url:  e.req.Url,
				})
				if err != nil {
					e.errc <- err
					continue
				}
				e.done <- data
			}
		}()
	}()
	return &Endpoint{channel: channel}
}

type Request struct {
	Host string `json:"host"`
	Url  string `json:"url"`
}

type Response struct {
	Host  string `json:"host"`
	Count int    `json:"count"`
}

func (e *Endpoint) check(w http.ResponseWriter, r *http.Request) {
	var req Request
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

	done := make(chan domain.StatsItem, 1)
	errc := make(chan error)
	item := channelItem{
		req:  req,
		done: done,
		errc: errc,
	}
	e.channel <- item
	select {
	case data := <-done:
		resp := Response{
			Host:  data.Host,
			Count: data.Count,
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			return
		}
	case <-errc:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
