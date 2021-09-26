package cache

import (
	"sync"

	"go.uber.org/zap"
)

type MemoryCache struct {
	sm  sync.Map
	log *zap.Logger
}

func NewMemoryCache(log *zap.Logger) *MemoryCache {
	return &MemoryCache{log: log}
}

func (m *MemoryCache) Get(host string) (int, bool) {
	v, ok := m.sm.Load(host)
	if ok {
		return v.(int), true
	}
	return 0, false
}

func (m *MemoryCache) Put(host string, n int) {
	m.log.Info("put", zap.String("host", host), zap.Int("count", n))
	m.sm.Store(host, n)
}
