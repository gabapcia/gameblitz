package rest

import (
	"net/http"

	"github.com/gabarcia/metagaming-api/internal/statistic"

	"github.com/gofiber/fiber/v2"
)

type UpsertPlayerStatisticReq struct {
	Value float64 `json:"value"` // Value that will be used to update the player's statistic
}

// @summary Upsert Player Statistic
// @description Set or update a player's statistic
// @router /api/v1/statistics/{statisticId}/players/{playerId} [POST]
// @accept json
// @produce json
// @param X-Game-ID header string true "Game ID responsible for the statistic"
// @param statisticId path string true "Statistic ID"
// @param playerId path string true "Player ID"
// @param UpsertPlayerRankData body UpsertPlayerRankReq true "Values to update the player rank"
// @success 204
// @failure 400,404,422,500 {object} ErrorResponse
func buildUpsertPlayerStatisticHandler(upsertPlayerStatisticFunc statistic.UpsertPlayerProgressionFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			statistic = c.Locals("statistic").(statistic.Statistic)
			playerID  = c.Params("playerId")
		)

		var body UpsertPlayerRankReq
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		if err := upsertPlayerStatisticFunc(c.Context(), statistic, playerID, body.Value); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	}
}
