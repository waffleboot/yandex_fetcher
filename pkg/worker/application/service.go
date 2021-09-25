package service

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/waffleboot/playstation_buy/pkg/common/domain"
)

type cache interface {
	Get(string) (int, bool)
	Put(string, int)
}

type initialService interface {
	Update(context.Context, string, int) error
}

type Service struct {
	clients        []http.Client
	initialService initialService
	cache          cache
}

func NewService(ctx context.Context, cache cache, initialService initialService, n int, timeout time.Duration) *Service {
	clients := make([]http.Client, n)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	for i := 0; i < n; i++ {
		clients[i].Timeout = timeout
		clients[i].Transport = tr
	}
	return &Service{clients: clients, cache: cache, initialService: initialService}
}

func (s *Service) Benchmark(item domain.YandexItem) (domain.StatsItem, error) {

	if n, ok := s.cache.Get(item.Host); ok {
		return domain.StatsItem{Host: item.Host, Count: n}, nil
	}

	req, err := http.NewRequest(http.MethodGet, item.Url, nil)
	if err != nil {
		return domain.StatsItem{}, err
	}

	var errCount uint32

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
				log.Printf("err %v", err)
				atomic.AddUint32(&errCount, 1)
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

	n := len(s.clients) - int(errCount)
	s.cache.Put(item.Host, n)
	s.initialService.Update(context.Background(), item.Host, n)

	return domain.StatsItem{
		Host:  item.Host,
		Count: n,
	}, nil
}
