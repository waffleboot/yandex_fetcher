package main

import (
	"log"
	"net"
	"sync"

	"github.com/valyala/fasthttp"
)

func main() {
	// N := 500
	N := 1
	var wg sync.WaitGroup
	wg.Add(N)
	var c fasthttp.Client
	c.MaxConnsPerHost = N
	ready := make(chan bool, N)
	start := make(chan bool, N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()
			req.SetRequestURI("https://tex-soyuz.ru/")
			c.Dial = func(addr string) (conn net.Conn, err error) {
				conn, err = fasthttp.Dial(addr)
				if err != nil {
					log.Println(err)
				}
				ready <- true
				<-start
				return
			}
			if err := c.Do(req, resp); err != nil {
				log.Println(err)
			}
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
		}()
	}
	for i := 0; i < N; i++ {
		<-ready
	}
	for i := 0; i < N; i++ {
		start <- true
	}
	wg.Wait()
}
