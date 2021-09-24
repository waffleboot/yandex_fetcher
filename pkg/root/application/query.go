package application

import (
	"context"
	"errors"
	"time"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

func (s *Service) Query(search string) (map[string]int, error) {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(s.timeout))
	data, err := s.supplier.GetYandexItems(ctx, search)
	if err != nil {
		return nil, err
	}
	m := make(map[string]int)
	p := make([]domain.YandexItem, 0, len(data))
	for _, v := range data {
		n := s.cache.Get(v.Host)
		if n > 0 {
			m[v.Host] = n
		} else {
			p = append(p, v)
		}
	}
	for _, v := range p {
		n, err := s.benchmark.Benchmark(ctx, v.Url)
		if err != nil {
			if !errors.Is(err, context.DeadlineExceeded) {
				return nil, err
			}
			break
		}
		s.cache.Put(v.Host, n)
		m[v.Host] = n
	}
	return m, nil
}
