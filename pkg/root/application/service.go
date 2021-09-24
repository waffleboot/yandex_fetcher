package application

import (
	"context"
	"time"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

type supplier interface {
	GetYandexItems(ctx context.Context, search string) ([]domain.YandexItem, error)
}

type benchmark interface {
	Benchmark(ctx context.Context, url string) (int, error)
}

type cache interface {
	Get(host string) int
	Put(host string, n int)
}

type Service struct {
	timeout   time.Duration
	benchmark benchmark
	supplier  supplier
	cache     cache
}

func NewService(
	timeout time.Duration,
	supplier supplier,
	benchmark benchmark,
	cache cache) *Service {
	return &Service{
		cache:     cache,
		timeout:   timeout,
		benchmark: benchmark,
		supplier:  supplier}
}
