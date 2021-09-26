package cache

type Cache interface {
	Get(host string) (int, bool)
	Put(host string, n int)
}
