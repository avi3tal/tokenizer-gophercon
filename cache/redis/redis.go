package redis

import (
	"context"
	"github.com/eran-levy/tokenizer-gophercon/cache"
	"github.com/eran-levy/tokenizer-gophercon/logger"
	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	client *redis.Client
	cfg    cache.Config
}

func New(cfg cache.Config) (cache.Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.CacheAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
		//ReadTimeout just an example, in real life you should set other options as well see redis.Options
		ReadTimeout: cfg.ReadTimeout,
		OnConnect: func(ctx context.Context, conn *redis.Conn) error {
			return conn.Ping(ctx).Err()
		},
	})
	err := rdb.Ping(context.Background()).Err()
	return redisCache{client: rdb, cfg: cfg}, err
}

func (r redisCache) Set(ctx context.Context, key string, value []byte) error {
	return r.client.Set(ctx, key, value, r.cfg.ExpirationTime).Err()
}

func (r redisCache) Get(ctx context.Context, key string) ([]byte, bool) {
	b, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return []byte{}, false
	} else if err != nil {
		logger.Log.With("cache_key", key).Errorf("could not retrieve from cache %s", err)
		return []byte{}, false
	} else {
		return b, true
	}
}

func (r redisCache) Close() error {
	return r.client.Close()
}
