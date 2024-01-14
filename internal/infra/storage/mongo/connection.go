package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type connection struct {
	client *mongo.Client
	db     string
}

func (c connection) ensureIndexes(ctx context.Context) error {
	if err := c.ensurePlayerStatisticIndexes(ctx); err != nil {
		return fmt.Errorf("Player Statistics: %w", err)
	}

	return nil
}

func (c connection) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func New(ctx context.Context, connStr, db string) (*connection, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.PrimaryPreferred()); err != nil {
		return nil, err
	}

	conn := &connection{
		client: client,
		db:     db,
	}
	return conn, conn.ensureIndexes(ctx)
}
