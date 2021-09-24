package application

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

type supplier = func(context.Context, string) ([]byte, error)

type Service struct {
	suppliers chan supplier
}

func NewService(factory func() supplier, n int) *Service {
	suppliers := make(chan supplier, n)
	for i := 0; i < n; i++ {
		suppliers <- factory()
	}
	return &Service{suppliers: suppliers}
}

func (s *Service) GetYandexItems(ctx context.Context, search string, done chan []domain.YandexItem, errorChannel chan error) {
	supplier := <-s.suppliers
	go func() {
		data, err := supplier(ctx, search)
		s.suppliers <- supplier
		if err != nil {
			errorChannel <- err
		}
		result := parseYandexResponse(data)
		if result.Error != nil {
			errorChannel <- result.Error
		}
		out := make([]domain.YandexItem, 0, len(result.Items))
		for _, v := range result.Items {
			out = append(out, domain.YandexItem{
				Host: v.Host,
				Url:  v.Url,
			})
		}
		done <- out
	}()
}
