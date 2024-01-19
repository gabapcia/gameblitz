package rest

import (
	"fmt"
	"time"

	"github.com/gabarcia/game-blitz/internal/leaderboard"
	"github.com/gabarcia/game-blitz/internal/quest"
	"github.com/gabarcia/game-blitz/internal/statistic"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	_ "github.com/gabarcia/game-blitz/internal/controller/rest/docs"
)

const (
	gameIDHeader = "X-Game-ID"
)

type Config struct {
	Port int

	CacheSorage               fiber.Storage
	CacheExpiration           time.Duration
	CacheMiddlewareExpiration time.Duration

	// Leaderboard
	CreateLeaderboardFunc              leaderboard.CreateFunc
	GetLeaderboardByIDAndGameIDFunc    leaderboard.GetByIDAndGameIDFunc
	DeleteLeaderboardByIDAndGameIDFunc leaderboard.SoftDeleteFunc

	UpsertPlayerRankFunc leaderboard.UpsertPlayerRankFunc
	RankingFunc          leaderboard.RankingFunc

	// Quest
	CreateQuestFunc           quest.CreateQuestFunc
	GetQuestByIDAndGameIDFunc quest.GetQuestByIDAndGameIDFunc
	SoftDeleteQuestFunc       quest.SoftDeleteQuestFunc

	StartQuestForPlayerFunc          quest.StartQuestForPlayerFunc
	GetPlayerQuestProgressionFunc    quest.GetPlayerQuestProgressionFunc
	UpdatePlayerQuestProgressionFunc quest.UpdatePlayerQuestProgressionFunc

	// Statistic
	CreateStatisticFunc                  statistic.CreateFunc
	GetStatisticByIDAndGameIDFunc        statistic.GetByIDAndGameIDFunc
	SoftDeleteStatisticByIDAndGameIDFunc statistic.SoftDeleteByIDAndGameIDFunc

	UpsertPlayerStatisticProgressionFunc statistic.UpsertPlayerProgressionFunc
	GetPlayerStatisticProgressionFunc    statistic.GetPlayerProgressionFunc
}

// @title Metagaming API
// @version 1.0
// @license.name MIT
// @description An API to handle basic gaming features like Statistics, Quests and Leaderboards
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

	rankings := leaderboards.Group("/:leaderboardId/ranking", buildGetLeaderboardMiddleware(config.CacheSorage, config.CacheMiddlewareExpiration, config.GetLeaderboardByIDAndGameIDFunc))
	rankings.Get("/", buildGetRankingHandler(config.RankingFunc))
	rankings.Post("/:playerId", buildUpsertPlayerRankHandler(config.UpsertPlayerRankFunc))

	// Quests
	quests := api.Group("/quests")
	quests.Post("/", buildCreateQuestHanlder(config.CreateQuestFunc))
	quests.Get("/:questId", buildGetQuestHanlder(config.GetQuestByIDAndGameIDFunc))
	quests.Delete("/:questId", buildDeleteQuestHanlder(config.SoftDeleteQuestFunc))

	playerQuests := quests.Group("/:questId/players", buildGetQuestMiddleware(config.CacheSorage, config.CacheMiddlewareExpiration, config.GetQuestByIDAndGameIDFunc))
	playerQuests.Post("/:playerId", buildStartPlayerQuestHandler(config.StartQuestForPlayerFunc))
	playerQuests.Get("/:playerId", buildGetPlayerQuestProgressionHandler(config.GetPlayerQuestProgressionFunc))
	playerQuests.Patch("/:playerId", buildUpdatePlayerQuestProgressionHandler(config.UpdatePlayerQuestProgressionFunc))

	// Statistic
	statistics := api.Group("/statistics")
	statistics.Post("/", buildCreateStatisticHandler(config.CreateStatisticFunc))
	statistics.Get("/:statisticId", buildGetStatisticHanlder(config.GetStatisticByIDAndGameIDFunc))
	statistics.Delete("/:statisticId", buildDeleteStatisticHanlder(config.SoftDeleteStatisticByIDAndGameIDFunc))

	playerStatistics := statistics.Group("/:statisticId/players", buildGetStatisticMiddleware(config.CacheSorage, config.CacheMiddlewareExpiration, config.GetStatisticByIDAndGameIDFunc))
	playerStatistics.Get("/:playerId", buildGetPlayerStatisticHandler(config.GetPlayerStatisticProgressionFunc))
	playerStatistics.Post("/:playerId", buildUpsertPlayerStatisticHandler(config.UpsertPlayerStatisticProgressionFunc))

	return app
}

func Execute(config Config) error {
	app := App(config)

	return app.Listen(fmt.Sprintf(":%d", config.Port))
}
