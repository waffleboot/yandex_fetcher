package service

import (
	"context"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Service struct {
	timeout time.Duration
	clients []http.Client
	ctx     context.Context
}

func NewService(ctx context.Context, n int, timeout time.Duration) *Service {
	clients := make([]http.Client, n)
	for i := 0; i < n; i++ {
		clients[i].Timeout = timeout
	}
	return &Service{clients: clients, ctx: ctx, timeout: timeout}
}

func (s *Service) Benchmark(url string) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	var counter uint32

	var wg sync.WaitGroup
	wg.Add(len(s.clients))

	start := make(chan bool, len(s.clients))

	for i := 0; i < len(s.clients); i++ {
		j := i
		go func() {
			defer wg.Done()
			<-start
			resp, err := s.clients[j].Do(req)
			if err != nil {
				atomic.AddUint32(&counter, 1)
			} else {
				defer resp.Body.Close()
			}
			log.Println(url)
		}()
	}

	for i := 0; i < len(s.clients); i++ {
		start <- true
	}

	wg.Wait()
	return int(counter), nil
}
