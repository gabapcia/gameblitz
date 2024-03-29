package main

import (
	"context"
	"time"

	"github.com/gabapcia/gameblitz/internal/auth"
	"github.com/gabapcia/gameblitz/internal/controller/rest"
	"github.com/gabapcia/gameblitz/internal/infra/async/rabbitmq"
	"github.com/gabapcia/gameblitz/internal/infra/cache/memcached"
	"github.com/gabapcia/gameblitz/internal/infra/logger/zap"
	"github.com/gabapcia/gameblitz/internal/infra/service/keycloack"
	"github.com/gabapcia/gameblitz/internal/infra/storage/mongo"
	"github.com/gabapcia/gameblitz/internal/infra/storage/postgres"
	"github.com/gabapcia/gameblitz/internal/infra/storage/redis"
	"github.com/gabapcia/gameblitz/internal/leaderboard"
	"github.com/gabapcia/gameblitz/internal/quest"
	"github.com/gabapcia/gameblitz/internal/statistic"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port int `envconfig:"PORT" required:"true"`

	KeycloackCertsURI string `envconfig:"KEYCLOACK_CERTS_URI" required:"true"`

	PotgresDSN string `envconfig:"POSTGRESQL_DSN" required:"true"`

	MongoURI string `envconfig:"MONGO_URI" required:"true"`
	MongoDB  string `envconfig:"MONGO_DB" required:"true"`

	RedisAddr     string `envconfig:"REDIS_ADDR" required:"true"`
	RedisUsername string `envconfig:"REDIS_USERNAME" required:"false"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" required:"false"`
	RedisDB       int    `envconfig:"REDIS_DB" required:"false"`

	MemcachedConnStr                   string `envconfig:"MEMCACHED_CONN_STR" required:"true"`
	MemcachedCacheExpiration           int    `envconfig:"MEMCACHED_EXPIRATION" required:"false" default:"60"`
	MemcachedCacheMiddlewareExpiration int    `envconfig:"MEMCACHED_MIDDLEWARE_EXPIRATION" required:"false" default:"60"`

	RabbitURI string `envconfig:"RABBITMQ_URI" required:"true"`
}

func main() {
	zap.Start()
	defer zap.Sync()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		zap.Panic(err, "env load failed")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	keycloack, err := keycloack.New(ctx, config.KeycloackCertsURI)
	if err != nil {
		zap.Panic(err, "keycloack startup failed")
	}

	redis := redis.New(ctx, config.RedisAddr, config.RedisUsername, config.RedisPassword, config.RedisDB)
	defer redis.Close()

	memcached := memcached.New(config.MemcachedConnStr)
	defer memcached.Close()

	rabbitmq, err := rabbitmq.NewProducer(ctx, config.RabbitURI)
	if err != nil {
		zap.Panic(err, "rabbitmq startup failed")
	}
	defer rabbitmq.Close()

	mongo, err := mongo.New(ctx, config.MongoURI, config.MongoDB)
	if err != nil {
		zap.Panic(err, "mongo startup failed")
	}
	defer mongo.Close(context.Background())

	postgres, err := postgres.New(ctx, config.PotgresDSN)
	if err != nil {
		zap.Panic(err, "postgres startup failed")
	}
	defer postgres.Close()

	restConfig := rest.Config{
		Port: config.Port,

		CacheSorage:               memcached,
		CacheExpiration:           time.Duration(config.MemcachedCacheExpiration) * time.Second,
		CacheMiddlewareExpiration: time.Duration(config.MemcachedCacheMiddlewareExpiration) * time.Second,

		// Auth
		AuthenticateFunc: auth.BuildAuthenticatorFunc(keycloack.Authenticate),

		// Leaderboard
		CreateLeaderboardFunc:              leaderboard.BuildCreateFunc(redis.CreateLeaderboard),
		GetLeaderboardByIDAndGameIDFunc:    leaderboard.BuildGetByIDAndGameIDFunc(redis.GetLeaderboardByIDAndGameID),
		DeleteLeaderboardByIDAndGameIDFunc: leaderboard.BuildSoftDeleteFunc(redis.SoftDeleteLeaderboard),

		UpsertPlayerRankFunc: leaderboard.BuildUpsertPlayerRankFunc(redis.UpsertPlayerRankValue),
		RankingFunc:          leaderboard.BuildRankingFunc(redis.GetRanking),

		// Quest
		CreateQuestFunc:           quest.BuildCreateQuestFunc(postgres.CreateQuest),
		GetQuestByIDAndGameIDFunc: quest.BuildGetQuestByIDAndGameIDFunc(postgres.GetQuestByIDAndGameID),
		SoftDeleteQuestFunc:       quest.BuildSoftDeleteQuestFunc(postgres.SoftDeleteQuestByIDAndGameID),

		StartQuestForPlayerFunc:          quest.BuildStartQuestForPlayerFunc(postgres.StartQuestForPlayer),
		GetPlayerQuestProgressionFunc:    quest.BuildGetPlayerQuestProgression(postgres.GetPlayerQuestProgression),
		UpdatePlayerQuestProgressionFunc: quest.BuildUpdatePlayerQuestProgressionFunc(rabbitmq.PlayerQuestProgressionUpdates, postgres.GetPlayerQuestProgression, postgres.UpdatePlayerQuestProgression),

		// Statistic
		CreateStatisticFunc:                  statistic.BuildCreateStatisticFunc(mongo.CreateStatistic),
		GetStatisticByIDAndGameIDFunc:        statistic.BuildGetStatisticByIDAndGameID(mongo.GetStatisticByIDAndGameID),
		SoftDeleteStatisticByIDAndGameIDFunc: statistic.BuildSoftDeleteStatistic(mongo.SoftDeleteStatistic),

		UpsertPlayerStatisticProgressionFunc: statistic.BuildUpsertPlayerProgressionFunc(rabbitmq.PlayerStatisticProgressionUpdates, mongo.UpdatePlayerStatisticProgression),
		GetPlayerStatisticProgressionFunc:    statistic.BuildGetPlayerProgression(mongo.GetPlayerProgression),
	}
	if err := rest.Execute(restConfig); err != nil {
		zap.Panic(err, "api execution failed")
	}
}
