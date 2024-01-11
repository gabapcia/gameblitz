package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type connection struct {
	rdb *redis.Client
}

func New(ctx context.Context, addr, username, password string, db int) *connection {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	conn := &connection{
		rdb: client,
	}

	return conn
}
