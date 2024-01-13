package memcached

import (
	"errors"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func (c connection) Get(key string) ([]byte, error) {
	data, err := c.client.Get(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			err = nil
		}

		return nil, err
	}

	return data.Value, nil
}

func (c connection) Set(key string, val []byte, exp time.Duration) error {
	return c.client.Set(&memcache.Item{
		Key:        key,
		Value:      val,
		Expiration: int32(exp.Seconds()),
	})
}

func (c connection) Delete(key string) error {
	if err := c.client.Delete(key); err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		return err
	}

	return nil
}

func (c connection) Reset() error {
	return c.client.DeleteAll()
}
