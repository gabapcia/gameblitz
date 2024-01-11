package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"
)

type Config struct {
	Port int

	CreateLeaderboardFunc              leaderboard.CreateFunc
	GetLeaderboardByIDAndGameIDFunc    leaderboard.GetByIDAndGameIDFunc
	DeleteLeaderboardByIDAndGameIDFunc leaderboard.SoftDeleteFunc
}

func Execute(config Config) error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	api := app.Group("/api/v1")

	leaderboards := api.Group("/leaderboards")
	leaderboards.Post("/", BuildCreateLeaderboardHandler(config.CreateLeaderboardFunc))
	leaderboards.Get("/:id", BuildGetLeaderboardHandler(config.GetLeaderboardByIDAndGameIDFunc))
	leaderboards.Delete("/:id", BuildDeleteLeaderboardHandler(config.DeleteLeaderboardByIDAndGameIDFunc))

	return app.Listen(fmt.Sprintf(":%d", config.Port))
}
