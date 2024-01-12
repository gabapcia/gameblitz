package rest

import (
	"fmt"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

const (
	gameIDHeader = "X-Game-ID"
)

type Config struct {
	Port   int
	Logger *zap.SugaredLogger

	CreateLeaderboardFunc              leaderboard.CreateFunc
	GetLeaderboardByIDAndGameIDFunc    leaderboard.GetByIDAndGameIDFunc
	DeleteLeaderboardByIDAndGameIDFunc leaderboard.SoftDeleteFunc
}

func App(config Config) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          BuildErrorHandler(config.Logger),
	})

	app.Use(recover.New())

	api := app.Group("/api/v1")

	leaderboards := api.Group("/leaderboards")
	leaderboards.Post("/", BuildCreateLeaderboardHandler(config.CreateLeaderboardFunc))
	leaderboards.Get("/:id", BuildGetLeaderboardHandler(config.GetLeaderboardByIDAndGameIDFunc))
	leaderboards.Delete("/:id", BuildDeleteLeaderboardHandler(config.DeleteLeaderboardByIDAndGameIDFunc))

	return app
}

func Execute(config Config) error {
	app := App(config)

	return app.Listen(fmt.Sprintf(":%d", config.Port))
}
