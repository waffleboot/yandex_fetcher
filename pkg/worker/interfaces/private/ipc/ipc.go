package ipc

import (
	"context"

	app "github.com/waffleboot/playstation_buy/pkg/worker/application"
)

type Endpoint struct {
	service *app.Service
}

func NewEndpoint(service *app.Service) *Endpoint {
	return &Endpoint{service: service}
}

func (e *Endpoint) Benchmark(ctx context.Context, url string) (int, error) {
	return e.service.Benchmark(ctx, url)
}
