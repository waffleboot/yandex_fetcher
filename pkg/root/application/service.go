package application

import (
	"context"
	"time"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

type yandex interface {
	GetItems(context.Context, string) ([]domain.YandexItem, error)
}

type benchmark interface {
	Benchmark(context.Context, []domain.YandexItem) (chan domain.StatsItem, chan error)
}

type cache interface {
	Get(string) int
	// Put(string, int)
}

type Service struct {
	timeout   time.Duration
	benchmark benchmark
	yandex    yandex
	cache     cache
}

func NewService(
	timeout time.Duration,
	yandex yandex,
	benchmark benchmark,
	cache cache) *Service {
	return &Service{
		cache:     cache,
		timeout:   timeout,
		benchmark: benchmark,
		yandex:    yandex}
}
