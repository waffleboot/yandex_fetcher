package memory

type MemoryCache struct {
	m map[string]int
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		m: make(map[string]int),
	}
}

func (m *MemoryCache) Get(host string) int {
	return m.m[host]
}

func (m *MemoryCache) Put(host string, n int) {
	m.m[host] = n
}
