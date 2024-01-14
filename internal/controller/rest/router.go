package rest

import (
	"fmt"
	"time"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"
	"github.com/gabarcia/metagaming-api/internal/quest"
	"github.com/gabarcia/metagaming-api/internal/ranking"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	_ "github.com/gabarcia/metagaming-api/internal/controller/rest/docs"
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

	CreateQuestFunc     quest.CreateQuestFunc
	SoftDeleteQuestFunc quest.SoftDeleteQuestFunc
}

// @title Metagaming API
// @version 1.0
// @license.name MIT
// @description An API to handle basic gaming features like Quests and Leaderboards
// @BasePath /
func App(config Config) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          buildErrorHandler(),
	})

	app.Use(recover.New())
	app.Get("/docs/*", swagger.HandlerDefault)
	app.Use(cache.New(cache.Config{
		Expiration:   config.CacheExpiration,
		Storage:      config.CacheSorage,
		CacheControl: true,
	}))

	api := app.Group("/api/v1")

	// Leaderboards
	leaderboards := api.Group("/leaderboards")
	leaderboards.Post("/", buildCreateLeaderboardHandler(config.CreateLeaderboardFunc))
	leaderboards.Get("/:leaderboardId", buildGetLeaderboardHandler(config.GetLeaderboardByIDAndGameIDFunc))
	leaderboards.Delete("/:leaderboardId", buildDeleteLeaderboardHandler(config.DeleteLeaderboardByIDAndGameIDFunc))

	rankings := leaderboards.Group(":leaderboardId/ranking", buildGetLeaderboardMiddleware(config.CacheSorage, config.CacheExpiration, config.GetLeaderboardByIDAndGameIDFunc))
	rankings.Get("/", buildGetRankingHandler(config.RankingFunc))
	rankings.Post("/:playerId", buildUpsertPlayerRankHandler(config.UpsertPlayerRankFunc))

	// Quests
	quests := api.Group("/quests")
	quests.Post("/", buildBuildCreateQuestHanlder(config.CreateQuestFunc))
	quests.Delete("/:questId", buildBuildDeleteQuestHanlder(config.SoftDeleteQuestFunc))

	return app
}

func Execute(config Config) error {
	app := App(config)

	return app.Listen(fmt.Sprintf(":%d", config.Port))
}
