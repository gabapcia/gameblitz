package rest

import (
	"net/http"

	"github.com/gabarcia/game-blitz/internal/leaderboard"

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

func rankFromDomain(r leaderboard.Rank) Rank {
	return Rank{
		PlayerID: r.PlayerID,
		Position: r.Position,
		Value:    r.Value,
	}
}

var (
	ErrorResponseLeaderboardClosed  = ErrorResponse{Code: "2.0", Message: "leaderboard closed"}
	ErrorResponseRankingPageNumber  = ErrorResponse{Code: "2.1", Message: "invalid page number"}
	ErrorResponseRankingLimitNumber = ErrorResponse{Code: "2.2", Message: "invalid limit number"}
)

// @summary Upsert Player Rank
// @description Set or update a player's rank on the leaderboard
// @router /api/v1/leaderboards/{leaderboardId}/ranking/{playerId} [POST]
// @accept json
// @produce json
// @param X-Game-ID header string true "Game ID responsible for the leaderboard"
// @param leaderboardId path string true "Leaderboard ID"
// @param playerId path string true "Player ID"
// @param UpsertPlayerRankData body UpsertPlayerRankReq true "Values to update the player rank"
// @success 204
// @failure 400,404,422,500 {object} ErrorResponse
func buildUpsertPlayerRankHandler(upsertPlayerRankFunc leaderboard.UpsertPlayerRankFunc) fiber.Handler {
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

// @summary Leaderboard Ranking
// @description Get the leaderboard ranking paginated
// @router /api/v1/leaderboards/{leaderboardId}/ranking [GET]
// @produce json
// @param X-Game-ID header string true "Game ID responsible for the leaderboard"
// @param leaderboardId path string true "Leaderboard ID"
// @param page query int false "Page number" minimun(0) default(0)
// @param limit query int false "Number of rankings per page" minimun(1) maximum(500) default(10)
// @success 200 {array} Rank
// @failure 400,404,422,500 {object} ErrorResponse
func buildGetRankingHandler(rankingFunc leaderboard.RankingFunc) fiber.Handler {
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
