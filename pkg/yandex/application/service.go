package application

import (
	"github.com/waffleboot/yandex_fetcher/pkg/common/domain"
)

type fetcher = func(string) ([]byte, error)

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

func (s *Service) GetItems(search string) ([]domain.YandexItem, error) {
	fetcher := <-s.fetchers
	data, err := fetcher(search)
	s.fetchers <- fetcher
	if err != nil {
		return nil, err
	}
	result := parseYandexResponse(data)
	if result.Error != nil {
		return nil, result.Error
	}
	out := make([]domain.YandexItem, 0, len(result.Items))
	for _, v := range result.Items {
		out = append(out, domain.YandexItem{
			Host: v.Host,
			Url:  v.Url,
		})
	}
	return out, nil
}
