package cache

import (
	"reflect"
	"time"

	"github.com/patrickmn/go-cache"
)

type MemoryStore struct {
	client            *cache.Cache
	DefaultExpiration time.Duration
}

type MemoryStoreOptions struct {
	DefaultExpiration time.Duration
	DefaultCacheItems map[string]cache.Item
	CleanupInterval   time.Duration
}

var MemoryStoreOptionsDefault = &MemoryStoreOptions{
	DefaultExpiration: time.Hour * 24,
	CleanupInterval:   time.Hour * 26,
}

func NewMemoryStore(options MemoryStoreOptions) *MemoryStore {
	var items = make(map[string]cache.Item)
	if options.DefaultCacheItems != nil {
		items = options.DefaultCacheItems
	}

	var client = cache.NewFrom(options.DefaultExpiration, options.CleanupInterval, items)
	return &MemoryStore{
		client:            client,
		DefaultExpiration: options.DefaultExpiration,
	}
}

func (c *MemoryStore) Get(key string, value interface{}) error {
	if !isPtr(value) {
		return ErrMustBePointer
	}

	val, found := c.client.Get(key)
	if !found {
		return ErrKeyNotFound
	}

	var i = reflect.ValueOf(val)
	var o = reflect.ValueOf(value)

	if i.Kind() != reflect.Ptr {
		i = toPtr(i)
	}

	if o.Kind() != reflect.Ptr {
		o = toPtr(o)
	}

	o.Elem().Set(i.Elem())

	return nil
}

func (c *MemoryStore) Set(key string, value interface{}, expiration ...time.Duration) error {
	if !isPtr(value) {
		return ErrMustBePointer
	}

	var exp = c.DefaultExpiration
	if len(expiration) > 0 {
		exp = expiration[0]
	}

	c.client.Set(key, value, exp)
	return nil
}

func (c *MemoryStore) Delete(key string) error {
	c.client.Delete(key)
	return nil
}

func (c *MemoryStore) Type() string {
	return "memory"
}
