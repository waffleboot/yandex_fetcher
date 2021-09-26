package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisCache struct {
	rdb *redis.Client
	log *zap.Logger
}

func NewRedisCache(redisAddr string, log *zap.Logger) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &RedisCache{log: log, rdb: rdb}
}

func (m *RedisCache) Get(host string) (int, bool) {
	v, err := m.rdb.Get(context.Background(), host).Int()
	if err != nil {
		return 0, false
	}
	return v, true
}

func (m *RedisCache) Put(host string, n int) {
	m.log.Info("put", zap.String("host", host), zap.Int("count", n))
	m.rdb.Set(context.Background(), host, n, 0)
}
