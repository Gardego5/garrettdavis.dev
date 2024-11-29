package bimarshal

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	Cache[T any] interface {
		Set(ctx context.Context, key string, data T, ttl time.Duration) error
		Get(ctx context.Context, key string) (*T, error)
		GetOrSet(ctx context.Context, key string, f func() (*T, time.Duration, error)) (*T, error)
	}
	redisCache[T any] struct {
		rdb *redis.Client
		pre string
		enc func(data *T) Bimarshal
	}
)

func NewCache[T any](r *redis.Client, pre string, enc func(data *T) Bimarshal) Cache[T] {
	return &redisCache[T]{r, pre, enc}
}

func (r *redisCache[T]) key(key string) string { return fmt.Sprintf("%s:%s", r.pre, key) }

func (r *redisCache[T]) Set(ctx context.Context, key string, data T, ttl time.Duration) error {
	return r.rdb.Set(ctx, r.key(key), r.enc(&data), ttl).Err()
}

func (r *redisCache[T]) Get(ctx context.Context, key string) (*T, error) {
	encoded, err := r.rdb.Get(ctx, r.key(key)).Result()
	if err != nil {
		return nil, err
	}

	data := new(T)
	err = r.enc(data).UnmarshalBinary([]byte(encoded))
	return data, err
}

func (r *redisCache[T]) GetOrSet(ctx context.Context, key string, f func() (*T, time.Duration, error)) (*T, error) {
	encoded, err := r.rdb.Get(ctx, r.key(key)).Result()
	if err == nil {
		data := new(T)
		err = r.enc(data).UnmarshalBinary([]byte(encoded))
		return data, err
	}
	data, ttl, err := f()
	if err != nil {
		return nil, err
	}
	if err = r.rdb.Set(ctx, r.key(key), r.enc(data), ttl).Err(); err != nil {
		return nil, err
	}
	return data, nil
}
