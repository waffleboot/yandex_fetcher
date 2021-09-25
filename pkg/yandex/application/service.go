package application

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

type parser = func(context.Context, string) ([]byte, error)

type Service struct {
	parsers chan parser
}

func NewService(factory func() parser, n int) *Service {
	parsers := make(chan parser, n)
	for i := 0; i < n; i++ {
		parsers <- factory()
	}
	return &Service{parsers: parsers}
}

func (s *Service) GetItems(ctx context.Context, search string, done chan []domain.YandexItem, errc chan error) {
	parser := <-s.parsers
	go func() {
		data, err := parser(ctx, search)
		s.parsers <- parser
		if err != nil {
			errc <- err
		}
		result := parseYandexResponse(data)
		if result.Error != nil {
			errc <- result.Error
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
