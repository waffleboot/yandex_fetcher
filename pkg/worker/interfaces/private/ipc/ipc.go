package ipc

import (
	app "github.com/waffleboot/playstation_buy/pkg/worker/application"
)

type Endpoint struct {
	service app.Service
}

func NewEndpoint(service app.Service) Endpoint {
	return Endpoint{service: service}
}

func (e Endpoint) Benchmark(url string) (int, error) {
	return e.service.Benchmark(url)
}
