package main

import (
	"sync"
	"testing"
)

func BenchmarkCond(b *testing.B) {

	var wg sync.WaitGroup
	wg.Add(b.N)

	ready := make(chan bool, b.N)

	v := sync.NewCond(&sync.RWMutex{})

	for i := 0; i < b.N; i++ {
		go func() {
			ready <- true
			v.L.Lock()
			v.Wait()
			v.L.Unlock()
			// op
			wg.Done()
		}()
	}
	for i := 0; i < b.N; i++ {
		<-ready
	}
	b.ResetTimer()
	v.Broadcast()
	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {

	var wg sync.WaitGroup
	wg.Add(b.N)

	ready := make(chan bool, b.N)
	start := make(chan bool, b.N)

	for i := 0; i < b.N; i++ {
		go func() {
			ready <- true
			<-start
			// op
			wg.Done()
		}()
	}
	for i := 0; i < b.N; i++ {
		<-ready
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start <- true
	}
	wg.Wait()
}
