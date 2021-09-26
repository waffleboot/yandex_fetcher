package application

import (
	"time"

	"github.com/waffleboot/yandex_fetcher/pkg/common/domain"
)

type yandex = func(string) ([]domain.YandexItem, error)

type benchmark = func(string, string) (int, error)

type cache interface {
	Get(string) (int, bool)
	Put(string, int)
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
