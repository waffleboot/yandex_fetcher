package http

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	"github.com/waffleboot/playstation_buy/pkg/worker/interfaces/private/http"
)

type BenchmarkSupplier struct {
	endpoint *http.Endpoint
}

func NewBenchmarkSupplier(endpoint *http.Endpoint) *BenchmarkSupplier {
	return &BenchmarkSupplier{
		endpoint: endpoint,
	}
}

func (b *BenchmarkSupplier) Benchmark(ctx context.Context, items []domain.YandexItem) (chan domain.StatsItem, chan error) {
	return b.endpoint.Benchmark(ctx, items)
}
