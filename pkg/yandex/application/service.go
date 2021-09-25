package application

import (
	"context"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

type fetcher = func(context.Context, string) ([]byte, error)

type Service struct {
	fetchers chan fetcher
}

func NewService(factory func() fetcher, n int) *Service {
	fetchers := make(chan fetcher, n)
	for i := 0; i < n; i++ {
		fetchers <- factory()
	}
	return &Service{fetchers: fetchers}
}

func (s *Service) GetItems(ctx context.Context, search string, done chan []domain.YandexItem, errc chan error) {
	fetcher := <-s.fetchers
	go func() {
		data, err := fetcher(ctx, search)
		s.fetchers <- fetcher
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
