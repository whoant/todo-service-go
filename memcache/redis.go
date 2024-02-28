package memcache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"todo-service/common"
	goservice "todo-service/go-sdk"
)

type redisCaching struct {
	store *cache.Cache
}

func NewRedisCaching(sc goservice.ServiceContext) *redisCaching {
	rdClient := sc.MustGet(common.PluginRedis).(*redis.Client)
	return &redisCaching{
		store: cache.New(&cache.Options{
			Redis:      rdClient,
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		}),
	}
}

func (r *redisCaching) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.store.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}

func (r *redisCaching) Get(ctx context.Context, key string, value interface{}) error {
	return r.store.Get(ctx, key, value)
}

func (r *redisCaching) Delete(ctx context.Context, key string) error {
	return r.store.Delete(ctx, key)
}
