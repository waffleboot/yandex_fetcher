package application

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

func (s *Service) fetchYandex(ctx context.Context, search string) ([]domain.YandexItem, error) {
	done := make(chan []domain.YandexItem, 1)
	errc := make(chan error, 1)
	go func() {
		data, err := s.yandex.GetItems(search)
		if err != nil {
			errc <- err
			return
		}
		done <- data
	}()
	select {
	case data := <-done:
		return data, nil
	case err := <-errc:
		return nil, fmt.Errorf("unable to fetch yandex page: %w", err)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *Service) ProcessQuery(search string) (map[string]int, error) {

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(s.timeout))
	defer cancel()

	data, err := s.fetchYandex(ctx, search)
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

	channel := make(chan domain.StatsItem, len(p))
	go func() {
		for _, v := range p {
			if count, ok := s.cache.Get(v.Host); ok {
				channel <- domain.StatsItem{
					Host:  v.Host,
					Count: count,
				}
				continue
			}
			count, err := s.benchmark.Benchmark(v.Host, v.Url)
			if err != nil {
				continue
			}
			s.cache.Put(v.Host, count)
			channel <- domain.StatsItem{
				Host:  v.Host,
				Count: count,
			}
		}
		close(channel)
	}()
	for {
		select {
		case item, ok := <-channel:
			if !ok {
				return m, nil
			}
			m[item.Host] = item.Count
		case <-ctx.Done():
			return m, ctx.Err()
		}
	}
}

func (s *Service) YandexItems(search string) ([]string, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(s.timeout))
	defer cancel()
	data, err := s.fetchYandex(ctx, search)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(data))
	for _, v := range data {
		out = append(out, v.Host)
	}
	sort.Strings(out)
	return out, nil
}

func (s *Service) CacheUpdate(host string, count int) {
	s.cache.Put(host, count)
}
