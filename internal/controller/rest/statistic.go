package rest

import (
	"net/http"
	"time"

	"github.com/gabarcia/metagaming-api/internal/statistic"
	"github.com/gofiber/fiber/v2"
)

type CreateStatisticReq struct {
	Name            string    `json:"name"`                                    // Statistic name
	Description     string    `json:"description"`                             // Statistic details
	AggregationMode string    `json:"aggregationMode" enums:"SUM,SUB,MAX,MIN"` // Data aggregation mode
	CanOverflow     bool      `json:"canOverflow"`                             // Can overflow the goal?
	Goal            *float64  `json:"goal"`                                    // Goal value. nil means no goal
	Landmarks       []float64 `json:"landmarks"`                               // Statistic landmarks
}

type Statistic struct {
	CreatedAt       time.Time `json:"createdAt"`                               // Time that the statistic was created
	UpdatedAt       time.Time `json:"updatedAt"`                               // Last time that the statistic was updated
	ID              string    `json:"id"`                                      // Statistic ID
	GameID          string    `json:"gameId"`                                  // ID of the game responsible for the statistic
	Name            string    `json:"name"`                                    // Statistic name
	Description     string    `json:"description"`                             // Statistic details
	AggregationMode string    `json:"aggregationMode" enums:"SUM,SUB,MAX,MIN"` // Data aggregation mode
	CanOverflow     bool      `json:"canOverflow"`                             // Can overflow the goal?
	Goal            *float64  `json:"goal"`                                    // Goal value. nil means no goal
	Landmarks       []float64 `json:"landmarks"`                               // Statistic landmarks
}

func (s CreateStatisticReq) toDomain(gameID string) statistic.NewStatisticData {
	return statistic.NewStatisticData{
		GameID:          gameID,
		Name:            s.Name,
		Description:     s.Description,
		AggregationMode: s.AggregationMode,
		CanOverflow:     s.CanOverflow,
		Goal:            s.Goal,
		Landmarks:       s.Landmarks,
	}
}

func statisticFromDomain(s statistic.Statistic) Statistic {
	return Statistic{
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		ID:              s.ID,
		GameID:          s.GameID,
		Name:            s.Name,
		Description:     s.Description,
		AggregationMode: s.AggregationMode,
		CanOverflow:     s.CanOverflow,
		Goal:            s.Goal,
		Landmarks:       s.Landmarks,
	}
}

var (
	ErrorResponseStatisticInvalid       = ErrorResponse{Code: "4.0", Message: "Invalid statistic"}
	ErrorResponseStatisticInvalidGameID = ErrorResponse{Code: "4.1", Message: "Invalid game id"}
	ErrorResponseStatisticNotFound      = ErrorResponse{Code: "4.2", Message: "Statistic not found"}
	ErrorResponseStatisticInvalidID     = ErrorResponse{Code: "4.3", Message: "Invalid statistic id"}
)

// @summary Create Statistic
// @description Create a statistic
// @router /api/v1/statistics [POST]
// @accept json
// @produce json
// @param X-Game-ID header string true "Game ID responsible for the leaderboard"
// @param NewStatisticData body CreateStatisticReq true "New statistic config data"
// @success 201 {object} Statistic
// @failure 400,422,500 {object} ErrorResponse
func buildCreateStatisticHandler(createStatisticFunc statistic.CreateFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		gameID := string(c.Request().Header.Peek(gameIDHeader))
		if gameID == "" {
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseStatisticInvalidGameID)
		}

		var body CreateStatisticReq
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		statistic, err := createStatisticFunc(c.Context(), body.toDomain(gameID))
		if err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(statisticFromDomain(statistic))
	}
}

// @summary Get Statistic By ID
// @description Get a statistic by its id
// @router /api/v1/statistics/{statisticId} [GET]
// @produce json
// @param X-Game-ID header string true "Game ID responsible for the leaderboard"
// @param statisticId path string true "Statistic ID"
// @success 200 {object} Statistic
// @failure 404,422,500 {object} ErrorResponse
func buildGetStatisticHanlder(getStatisticByIDAndGameID statistic.GetByIDAndGameID) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			statisticID = c.Params("statisticId")
			gameID      = string(c.Request().Header.Peek(gameIDHeader))
		)

		if gameID == "" {
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseStatisticInvalidGameID)
		}

		statistic, err := getStatisticByIDAndGameID(c.Context(), statisticID, gameID)
		if err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(statisticFromDomain(statistic))
	}
}

// @summary Delete Statistic
// @description Delete a statistic by its id
// @router /api/v1/statistics/{statisticId} [DELETE]
// @param X-Game-ID header string true "Game ID responsible for the leaderboard"
// @param statisticId path string true "Statistic ID"
// @success 204
// @failure 404,422,500 {object} ErrorResponse
func buildDeleteStatisticHanlder(softDeleteStatisticFunc statistic.SoftDeleteByIDAndGameID) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			questID = c.Params("statisticId")
			gameID  = string(c.Request().Header.Peek(gameIDHeader))
		)

		if gameID == "" {
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseStatisticInvalidGameID)
		}

		if err := softDeleteStatisticFunc(c.Context(), questID, gameID); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	}
}
