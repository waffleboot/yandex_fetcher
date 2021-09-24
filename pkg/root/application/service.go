package application

import (
	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

type supplier interface {
	Supply(search string) ([]domain.YandexItem, error)
}

type benchmark interface {
	Benchmark(url string) (int, error)
}

type Service struct {
	benchmark benchmark
	supplier  supplier
}

func NewService(supplier supplier, benchmark benchmark) *Service {
	return &Service{
		benchmark: benchmark,
		supplier:  supplier}
}
