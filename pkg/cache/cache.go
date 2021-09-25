package cache

import "sync"

type MemoryCache struct {
	sm sync.Map
}

func (m *MemoryCache) Get(host string) int {
	v, ok := m.sm.Load(host)
	if ok {
		return v.(int)
	}
	return 0
}

func (m *MemoryCache) Put(host string, n int) {
	m.sm.Store(host, n)
}
