package main

import (
	"sync"
	"testing"
)

func BenchmarkCond(b *testing.B) {

	var wg sync.WaitGroup
	wg.Add(b.N)

	var start bool
	ready := make(chan bool)
	v := sync.NewCond(&sync.Mutex{})

	for i := 0; i < b.N; i++ {
		go func() {
			ready <- true
			v.L.Lock()
			if !start {
				v.Wait()
			}
			v.L.Unlock()
			// op
			wg.Done()
		}()
	}
	for i := 0; i < b.N; i++ {
		<-ready
	}
	start = true
	b.ResetTimer()
	v.Broadcast()
	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {

	var wg sync.WaitGroup
	wg.Add(b.N)

	ready := make(chan bool, b.N)
	start := make(chan bool)

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
	close(start)
	wg.Wait()
}

func BenchmarkLock(b *testing.B) {

	var wg sync.WaitGroup
	wg.Add(b.N)

	ready := make(chan bool, b.N)

	var l sync.RWMutex
	l.Lock()

	for i := 0; i < b.N; i++ {
		go func() {
			ready <- true
			l.RLock()
			// op
			wg.Done()
		}()
	}
	for i := 0; i < b.N; i++ {
		<-ready
	}
	b.ResetTimer()
	l.Unlock()
	wg.Wait()
}
