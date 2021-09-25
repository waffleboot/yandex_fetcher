package cache

import (
	"log"
	"sync"
)

type MemoryCache struct {
	sm sync.Map
}

func (m *MemoryCache) Get(host string) (int, bool) {
	v, ok := m.sm.Load(host)
	if ok {
		return v.(int), true
	}
	return 0, false
}

func (m *MemoryCache) Put(host string, n int) {
	log.Printf("put %s %d", host, n)
	m.sm.Store(host, n)
}
