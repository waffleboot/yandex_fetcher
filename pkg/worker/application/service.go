package service

import (
	"context"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

type cache interface {
	Put(string, int)
}

type Service struct {
	clients []http.Client
	cache   cache
}

func NewService(ctx context.Context, cache cache, n int, timeout time.Duration) *Service {
	clients := make([]http.Client, n)
	for i := 0; i < n; i++ {
		clients[i].Timeout = timeout
	}
	return &Service{clients: clients, cache: cache}
}

func (s *Service) Benchmark(item domain.YandexItem) (domain.StatsItem, error) {

	log.Printf("test %s", item.Host)
	req, err := http.NewRequest(http.MethodGet, item.Url, nil)
	if err != nil {
		return domain.StatsItem{}, err
	}

	var counter uint32

	var wg sync.WaitGroup
	wg.Add(len(s.clients))

	ready := make(chan bool, len(s.clients))
	start := make(chan bool, len(s.clients))

	for i := 0; i < len(s.clients); i++ {
		j := i
		go func() {
			defer wg.Done()
			ready <- true
			<-start
			resp, err := s.clients[j].Do(req)
			if err != nil {
				atomic.AddUint32(&counter, 1)
			} else {
				defer resp.Body.Close()
			}
		}()
	}
	for i := 0; i < len(s.clients); i++ {
		<-ready
	}
	for i := 0; i < len(s.clients); i++ {
		start <- true
	}
	wg.Wait()
	n := len(s.clients) - int(counter)
	s.cache.Put(item.Host, n)
	return domain.StatsItem{
		Host:  item.Host,
		Count: n,
	}, nil
}
