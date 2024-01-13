package memcached

import "time"

func (c connection) Get(key string) ([]byte, error) {
	return c.Get(key)
}

func (c connection) Set(key string, val []byte, exp time.Duration) error {
	return c.Set(key, val, exp)
}

func (c connection) Delete(key string) error {
	return c.Delete(key)
}

func (c connection) Reset() error {
	return c.client.DeleteAll()
}
