package application

import (
	"context"
	"errors"
	"time"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

var ErrInvalidChecker = errors.New("invalid checker")

func (s *Service) ProcessQuery(search string) (map[string]int, error) {

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(s.timeout))
	defer cancel()

	data, err := s.yandex.GetItems(ctx, search)
	if err != nil {
		return nil, err
	}

	m := make(map[string]int)
	p := make([]domain.YandexItem, 0, len(data))
	for _, v := range data {
		if n, ok := s.cache.Get(v.Host); ok {
			m[v.Host] = n
		} else {
			p = append(p, v)
		}
	}
	if len(p) == 0 {
		return m, nil
	}

	items, err := s.benchmark.Benchmark(ctx, p)
	for _, v := range items {
		s.cache.Put(v.Host, v.Count)
		m[v.Host] = v.Count
	}
	return m, err
}

func (s *Service) Update(host string, count int) {
	s.cache.Put(host, count)
}
