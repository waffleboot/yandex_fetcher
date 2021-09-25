package application

import (
	"context"
	"log"
	"time"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

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
		if n := s.cache.Get(v.Host); n > 0 {
			m[v.Host] = n
		} else {
			p = append(p, v)
		}
	}
	if len(p) == 0 {
		return m, nil
	}

	log.Printf("%d not found", len(p))
	datc, errc := s.benchmark.Benchmark(ctx, p)
	for {
		select {
		case d := <-datc:
			m[d.Host] = d.Count
		case err := <-errc:
			return m, err
		case <-ctx.Done():
			return m, ctx.Err()
		}
	}
}
