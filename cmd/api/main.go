package main

import (
	"context"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/gabarcia/metagaming-api/internal/controller/rest"
	"github.com/gabarcia/metagaming-api/internal/infra/cache/memcached"
	"github.com/gabarcia/metagaming-api/internal/infra/logger/zap"
	"github.com/gabarcia/metagaming-api/internal/infra/storage/postgres"
	"github.com/gabarcia/metagaming-api/internal/infra/storage/redis"
	"github.com/gabarcia/metagaming-api/internal/leaderboard"
	"github.com/gabarcia/metagaming-api/internal/ranking"
)

type Config struct {
	Port int `envconfig:"PORT" required:"true"`

	PotgresDSN string `envconfig:"POSTGRESQL_DSN" required:"true"`

	RedisAddr     string `envconfig:"REDIS_ADDR" required:"true"`
	RedisUsername string `envconfig:"REDIS_USERNAME" required:"false"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" required:"false"`
	RedisDB       int    `envconfig:"REDIS_DB" required:"false"`

	MemcachedConnStr         string `envconfig:"MEMCACHED_CONN_STR" required:"true"`
	MemcachedCacheExpiration int    `envconfig:"MEMCACHED_EXPIRATION" required:"true"`
}

func main() {
	zap.Start()
	defer zap.Sync()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		zap.Panic(err, "env load failed")
	}

	ctx := context.Background()

	redis := redis.New(ctx, config.RedisAddr, config.RedisUsername, config.RedisPassword, config.RedisDB)
	defer redis.Close()

	memcached := memcached.New(config.MemcachedConnStr)
	defer memcached.Close()

	postgres, err := postgres.New(ctx, config.PotgresDSN)
	if err != nil {
		zap.Panic(err, "postgres startup failed")
	}
	defer postgres.Close()

	restConfig := rest.Config{
		Port: config.Port,

		CacheSorage:     memcached,
		CacheExpiration: time.Duration(config.MemcachedCacheExpiration) * time.Second,

		CreateLeaderboardFunc:              leaderboard.BuildCreateFunc(postgres.CreateLeaderboard),
		GetLeaderboardByIDAndGameIDFunc:    leaderboard.BuildGetByIDAndGameIDFunc(postgres.GetLeaderboardByIDAndGameID),
		DeleteLeaderboardByIDAndGameIDFunc: leaderboard.BuildSoftDeleteFunc(postgres.SoftDeleteLeaderboard),

		UpsertPlayerRankFunc: ranking.BuildUpsertPlayerRankFunc(redis.IncrementPlayerRankValue, redis.SetMaxPlayerRankValue, redis.SetMinPlayerRankValue),
		RankingFunc:          ranking.BuildRankingFunc(redis.GetRanking),
	}
	if err := rest.Execute(restConfig); err != nil {
		zap.Panic(err, "api execution failed")
	}
}
