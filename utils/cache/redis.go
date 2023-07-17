package cache

import (
	"context"
	"fileCollect/global"
	"time"
)

type RedisStore struct {
	Expiration time.Duration
	Context	   context.Context
}

func (r *RedisStore) Set(key, value string) error {
	return global.Rdb.Set(r.Context, key, value, r.Expiration).Err()
}

func (r *RedisStore) Get(key string) (res string, err error) {
	res, err = global.Rdb.Get(r.Context, key).Result()
	return
}

func (r *RedisStore) Del(key string) error {
	return global.Rdb.Del(r.Context, key).Err()
}

// return redisStore structure
func SetRedisStore(ctx context.Context, expire time.Duration) *RedisStore {
	return &RedisStore{
		Expiration: expire,
		Context: ctx,
	}
}