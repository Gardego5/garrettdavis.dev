package bimarshal

import (
	"reflect"

	"github.com/redis/go-redis/v9"
)

type (
	registration[T any] struct{ enc func(data *T) Bimarshal }
	register            interface {
		register(*redis.Client, map[reflect.Type]any, string)
	}
	Caches           map[string]register
	RegisteredCaches map[reflect.Type]any
)

func (c Caches) Build(rdb *redis.Client) RegisteredCaches {
	m := make(RegisteredCaches)
	for k, v := range c {
		v.register(rdb, m, k)
	}
	return m
}

func (r registration[T]) register(rdb *redis.Client, c map[reflect.Type]any, prefix string) {
	c[reflect.TypeFor[T]()] = NewCache[T](rdb, prefix, r.enc)
}

func Register[T any](enc func(data *T) Bimarshal) registration[T] { return registration[T]{enc} }

func Get[T any](caches RegisteredCaches) Cache[T] { return caches[reflect.TypeFor[T]()].(Cache[T]) }
