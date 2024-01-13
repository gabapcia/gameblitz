package rest

import (
	"fmt"
	"time"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"
	"github.com/gabarcia/metagaming-api/internal/ranking"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

const (
	gameIDHeader = "X-Game-ID"
)

type Config struct {
	Port int

	CacheSorage     fiber.Storage
	CacheExpiration time.Duration

	CreateLeaderboardFunc              leaderboard.CreateFunc
	GetLeaderboardByIDAndGameIDFunc    leaderboard.GetByIDAndGameIDFunc
	DeleteLeaderboardByIDAndGameIDFunc leaderboard.SoftDeleteFunc

	UpsertPlayerRankFunc ranking.UpsertPlayerRankFunc
	RankingFunc          ranking.RankingFunc
}

func App(config Config) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          BuildErrorHandler(),
	})

	app.Use(recover.New())
	app.Use(cache.New(cache.Config{
		Expiration: config.CacheExpiration,
		Storage:    config.CacheSorage,
	}))

	api := app.Group("/api/v1")

	leaderboards := api.Group("/leaderboards")
	leaderboards.Post("/", BuildCreateLeaderboardHandler(config.CreateLeaderboardFunc))
	leaderboards.Get("/:leaderboardId", BuildGetLeaderboardHandler(config.GetLeaderboardByIDAndGameIDFunc))
	leaderboards.Delete("/:leaderboardId", BuildDeleteLeaderboardHandler(config.DeleteLeaderboardByIDAndGameIDFunc))

	rankings := leaderboards.Group(":leaderboardId/ranking", BuildGetLeaderboardMiddleware(config.CacheSorage, config.CacheExpiration, config.GetLeaderboardByIDAndGameIDFunc))
	rankings.Get("/", BuildGetRankingHandler(config.RankingFunc))
	rankings.Post("/:playerId", BuildUpsertPlayerRankHandler(config.UpsertPlayerRankFunc))

	return app
}

func Execute(config Config) error {
	app := App(config)

	return app.Listen(fmt.Sprintf(":%d", config.Port))
}
