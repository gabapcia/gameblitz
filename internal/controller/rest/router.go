package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v3"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"
)

type Config struct {
	Port int

	CreateLeaderboardFunc leaderboard.CreateFunc
}

func Execute(config Config) error {
	app := fiber.New()

	api := app.Group("/api/v1")

	leaderboards := api.Group("/leaderboards")
	leaderboards.Post("/", BuildCreateLeaderboardHandler(config.CreateLeaderboardFunc))

	return app.Listen(fmt.Sprintf(":%d", config.Port))
}