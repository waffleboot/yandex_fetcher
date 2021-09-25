package ipc

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
	"github.com/waffleboot/playstation_buy/pkg/worker/interfaces/private/ipc"
)

type BenchmarkSupplier struct {
	endpoint *ipc.Endpoint
}

func NewBenchmarkSupplier(endpoint *ipc.Endpoint) *BenchmarkSupplier {
	return &BenchmarkSupplier{
		endpoint: endpoint,
	}
}

func (b *BenchmarkSupplier) Benchmark(ctx context.Context, items []domain.YandexItem) (chan domain.StatsItem, chan error) {
	return b.endpoint.Benchmark(ctx, items)
}
