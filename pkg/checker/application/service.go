package service

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

type cache interface {
	Get(string) (int, bool)
	Put(string, int)
}

type initialService interface {
	CacheUpdate(string, int) error
}

type Service struct {
	clients        int
	client         fasthttp.Client
	initialService initialService
	cache          cache
	token          chan bool
}

func NewService(cache cache, initialService initialService, n int, timeout time.Duration) *Service {
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }

	// client.Timeout = timeout
	// client.Transport = tr

	return &Service{
		cache:   cache,
		clients: n,
		client: fasthttp.Client{
			MaxConnsPerHost: n,
		},
		token:          make(chan bool, 1),
		initialService: initialService}
}

func (s *Service) Benchmark(host, url string) (int, error) {

	if n, ok := s.cache.Get(host); ok {
		return n, nil
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethodBytes([]byte(http.MethodGet))
	req.Header.Set("Connection", "close")

	var errCount uint32

	var wg sync.WaitGroup
	wg.Add(s.clients)

	ready := make(chan bool, s.clients)
	start := make(chan bool, s.clients)

	s.token <- true
	if n, ok := s.cache.Get(host); ok {
		<-s.token
		return n, nil
	}
	for i := 0; i < s.clients; i++ {
		go func() {
			defer wg.Done()
			resp := fasthttp.AcquireResponse()
			ready <- true
			<-start

			err := s.client.Do(req, resp)
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)

			if err != nil {
				atomic.AddUint32(&errCount, 1)
			}
		}()
	}
	for i := 0; i < s.clients; i++ {
		<-ready
	}
	for i := 0; i < s.clients; i++ {
		start <- true
	}
	wg.Wait()
	n := s.clients - int(errCount)
	s.cache.Put(host, n)
	<-s.token

	s.initialService.CacheUpdate(host, n)

	return n, nil
}
