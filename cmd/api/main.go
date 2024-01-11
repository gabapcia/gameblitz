package main

import (
	"context"

	"github.com/kelseyhightower/envconfig"

	"github.com/gabarcia/metagaming-api/internal/controller/rest"
	"github.com/gabarcia/metagaming-api/internal/infra/storage/postgres"
	"github.com/gabarcia/metagaming-api/internal/leaderboard"
)

type Config struct {
	Port       int    `envconfig:"PORT" required:"true"`
	PotgresDSN string `envconfig:"POSTGRESQL_DSN" required:"true"`
}

func main() {
	ctx := context.Background()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}

	postgres, err := postgres.New(ctx, config.PotgresDSN)
	if err != nil {
		panic(err)
	}
	defer postgres.Close()

	restConfig := rest.Config{
		Port:                               config.Port,
		CreateLeaderboardFunc:              leaderboard.BuildCreateFunc(postgres.CreateLeaderboard),
		GetLeaderboardByIDAndGameIDFunc:    leaderboard.BuildGetByIDAndGameIDFunc(postgres.GetLeaderboardByIDAndGameID),
		DeleteLeaderboardByIDAndGameIDFunc: leaderboard.BuildSoftDeleteFunc(postgres.SoftDeleteLeaderboard),
	}
	if err := rest.Execute(restConfig); err != nil {
		panic(err)
	}
}
