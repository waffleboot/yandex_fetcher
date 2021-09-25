package service

import (
	"crypto/tls"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type cache interface {
	Get(string) (int, bool)
	Put(string, int)
}

type initialService interface {
	CacheUpdate(string, int) error
}

type Service struct {
	clients        []http.Client
	initialService initialService
	cache          cache
	token          chan bool
}

func NewService(cache cache, initialService initialService, n int, timeout time.Duration) *Service {
	clients := make([]http.Client, n)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	for i := 0; i < n; i++ {
		clients[i].Timeout = timeout
		clients[i].Transport = tr
	}
	return &Service{
		cache:          cache,
		clients:        clients,
		token:          make(chan bool, 1),
		initialService: initialService}
}

func (s *Service) Benchmark(host, url string) (int, error) {

	if n, ok := s.cache.Get(host); ok {
		return n, nil
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	var errCount uint32

	var wg sync.WaitGroup
	wg.Add(len(s.clients))

	ready := make(chan bool, len(s.clients))
	start := make(chan bool, len(s.clients))

	s.token <- true
	if n, ok := s.cache.Get(host); ok {
		<-s.token
		return n, nil
	}
	for i := 0; i < len(s.clients); i++ {
		j := i
		go func() {
			defer wg.Done()
			ready <- true
			<-start
			resp, err := s.clients[j].Do(req)
			if err != nil {
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
	s.cache.Put(host, n)
	<-s.token

	s.initialService.CacheUpdate(host, n)

	return n, nil
}
