package service

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type cache interface {
	Get(string) (int, bool)
	Put(string, int)
}

type initialService = func(string, int) error

type Service struct {
	clients        []http.Client
	initialService initialService
	cache          cache
	token          sync.RWMutex
}

func NewService(cache cache, initialService initialService, n int, timeout time.Duration) *Service {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	clients := make([]http.Client, n)
	for i := 0; i < n; i++ {
		clients[i].Timeout = 30 * time.Second
		clients[i].Transport = tr
		clients[i].Jar = &MyJar{}
	}

	return &Service{
		cache:          cache,
		clients:        clients,
		initialService: initialService}
}

func (s *Service) Benchmark(host, url string) (int, error) {

	if n, ok := s.cache.Get(host); ok {
		return n, nil
	}

	var errCount uint32

	var wg sync.WaitGroup
	wg.Add(len(s.clients))

	ready := make(chan bool, len(s.clients))
	start := make(chan bool, len(s.clients))

	s.token.Lock()
	if n, ok := s.cache.Get(host); ok {
		s.token.Unlock()
		return n, nil
	}
	for i := 0; i < len(s.clients); i++ {
		j := i
		go func() {
			defer wg.Done()

			req, err := http.NewRequest(http.MethodGet, url, nil)

			ready <- true

			if err != nil {
				log.Printf("error on %s %v", url, err)
				atomic.AddUint32(&errCount, 1)
				return
			}
			req.Header.Set("Connection", "close")
			req.Header.Set("User-Agent", "Mozilla/5.0 Gecko/20100101 Firefox/92.0")

			<-start

			resp, err := s.clients[j].Do(req)
			if err != nil {
				log.Printf("error on %s %v", url, err)
				atomic.AddUint32(&errCount, 1)
				return
			}
			defer resp.Body.Close()
		}()
	}
	for i := 0; i < len(s.clients); i++ {
		<-ready
	}
	for i := 0; i < len(s.clients); i++ {
		start <- true
	}
	wg.Wait()
	n := len(s.clients) - int(atomic.LoadUint32(&errCount))
	s.cache.Put(host, n)
	s.token.Unlock()

	s.initialService(host, n)

	return n, nil
}

type MyJar struct {
	cookies []*http.Cookie
}

func (j *MyJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.cookies = cookies
}

func (j *MyJar) Cookies(u *url.URL) []*http.Cookie {
	return j.cookies
}
