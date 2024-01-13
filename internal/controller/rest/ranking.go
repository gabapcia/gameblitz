package rest

import (
	"net/http"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"
	"github.com/gabarcia/metagaming-api/internal/ranking"

	"github.com/gofiber/fiber/v2"
)

type UpsertPlayerRankReq struct {
	Value float64 `json:"value"` // Value that will be used to update the player's rank
}

type Rank struct {
	PlayerID string  `json:"playerId"` // Player's ID
	Position int64   `json:"position"` // Player ranking position
	Value    float64 `json:"value"`    // Player rank value
}

func rankFromDomain(r ranking.Rank) Rank {
	return Rank{
		PlayerID: r.PlayerID,
		Position: r.Position,
		Value:    r.Value,
	}
}

var (
	ErrorResponseLeaderboardClosed = ErrorResponse{Code: "2.0", Message: "leaderboard closed"}
)

func buildUpsertPlayerRankHandler(upsertPlayerRankFunc ranking.UpsertPlayerRankFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			leaderboard = c.Locals("leaderboard").(leaderboard.Leaderboard)
			playerID    = c.Params("playerId")
		)

		var body UpsertPlayerRankReq
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		if err := upsertPlayerRankFunc(c.Context(), leaderboard, playerID, body.Value); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	}
}

func buildGetRankingHandler(rankingFunc ranking.RankingFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			leaderboard = c.Locals("leaderboard").(leaderboard.Leaderboard)
			page        = c.QueryInt("page", 0)
			limit       = c.QueryInt("limit", 10)
		)

		rankings, err := rankingFunc(c.Context(), leaderboard, int64(page), int64(limit))
		if err != nil {
			return err
		}

		data := make([]Rank, len(rankings))
		for i, rank := range rankings {
			data[i] = rankFromDomain(rank)
		}

		return c.Status(http.StatusOK).JSON(data)
	}
}
