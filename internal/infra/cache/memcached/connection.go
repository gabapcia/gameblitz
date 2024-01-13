package memcached

import "github.com/bradfitz/gomemcache/memcache"

type connection struct {
	client *memcache.Client
}

func (c connection) Close() error {
	return c.client.Close()
}

func New(connStr string) *connection {
	client := memcache.New(connStr)

	conn := &connection{
		client: client,
	}
	return conn
}
