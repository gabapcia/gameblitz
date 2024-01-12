package main

import (
	"context"
	"os"

	"github.com/kelseyhightower/envconfig"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"

	"github.com/gabarcia/metagaming-api/internal/controller/rest"
	"github.com/gabarcia/metagaming-api/internal/infra/storage/postgres"
	"github.com/gabarcia/metagaming-api/internal/leaderboard"
)

type Config struct {
	Port       int    `envconfig:"PORT" required:"true"`
	PotgresDSN string `envconfig:"POSTGRESQL_DSN" required:"true"`
}

func main() {
	var (
		ctx = context.Background()

		core   = ecszap.NewCore(ecszap.NewDefaultEncoderConfig(), os.Stdout, zap.InfoLevel)
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.PanicLevel)).Sugar()
	)
	defer logger.Sync()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		logger.Panic(err)
	}

	postgres, err := postgres.New(ctx, config.PotgresDSN)
	if err != nil {
		logger.Panic(err)
	}
	defer postgres.Close()

	restConfig := rest.Config{
		Port:                               config.Port,
		CreateLeaderboardFunc:              leaderboard.BuildCreateFunc(postgres.CreateLeaderboard),
		GetLeaderboardByIDAndGameIDFunc:    leaderboard.BuildGetByIDAndGameIDFunc(postgres.GetLeaderboardByIDAndGameID),
		DeleteLeaderboardByIDAndGameIDFunc: leaderboard.BuildSoftDeleteFunc(postgres.SoftDeleteLeaderboard),
	}
	if err := rest.Execute(restConfig); err != nil {
		logger.Panic(err)
	}
}
