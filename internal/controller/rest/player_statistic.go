package rest

import (
	"net/http"
	"time"

	"github.com/gabarcia/metagaming-api/internal/statistic"

	"github.com/gofiber/fiber/v2"
)

type UpsertPlayerStatisticReq struct {
	Value float64 `json:"value"` // Value that will be used to update the player's statistic
}

type (
	PlayerStatisticProgressionLandmark struct {
		Value       float64    `json:"value"`                 // Landmark value
		Completed   bool       `json:"completed"`             // Has the player reached the landmark?
		CompletedAt *time.Time `json:"completedAt,omitempty"` // Time the player reached the landmark
	}

	PlayerStatisticProgression struct {
		StartedAt       *time.Time                           `json:"startedAt,omitempty"`       // Time the player started the progression for the given statistic
		UpdatedAt       *time.Time                           `json:"updatedAt,omitempty"`       // Last time the player updated it's statistic progress
		PlayerID        string                               `json:"playerId"`                  // Player's ID
		StatisticID     string                               `json:"statisticId"`               // Statistic ID
		CurrentValue    *float64                             `json:"currentValue"`              // Current progression value
		GoalValue       *float64                             `json:"goalValue"`                 // Statistic's goal
		GoalCompleted   *bool                                `json:"goalCompleted,omitempty"`   // Has the player reached the goal?
		GoalCompletedAt *time.Time                           `json:"goalCompletedAt,omitempty"` // Time the player reached the goal
		Landmarks       []PlayerStatisticProgressionLandmark `json:"landmarks"`                 // Landmarks player progression
	}
)

func playerStatisticProgressionFromDomain(p statistic.PlayerProgression) PlayerStatisticProgression {
	landmarks := make([]PlayerStatisticProgressionLandmark, len(p.Landmarks))
	for i, landmark := range p.Landmarks {
		var completedAt *time.Time
		if !landmark.CompletedAt.IsZero() {
			tmp := landmark.CompletedAt
			completedAt = &tmp
		}

		landmarks[i] = PlayerStatisticProgressionLandmark{
			Value:       landmark.Value,
			Completed:   landmark.Completed,
			CompletedAt: completedAt,
		}
	}

	var startedAt *time.Time
	if !p.StartedAt.IsZero() {
		startedAt = &p.StartedAt
	}

	var updatedAt *time.Time
	if !p.UpdatedAt.IsZero() {
		updatedAt = &p.UpdatedAt
	}

	var goalCompletedAt *time.Time
	if !p.GoalCompletedAt.IsZero() {
		goalCompletedAt = &p.GoalCompletedAt
	}

	return PlayerStatisticProgression{
		StartedAt:       startedAt,
		UpdatedAt:       updatedAt,
		PlayerID:        p.PlayerID,
		StatisticID:     p.StatisticID,
		CurrentValue:    p.CurrentValue,
		GoalValue:       p.GoalValue,
		GoalCompleted:   p.GoalCompleted,
		GoalCompletedAt: goalCompletedAt,
		Landmarks:       landmarks,
	}
}

var (
	ErrorResponsePlayerStatisticNotFound = ErrorResponse{Code: "5.0", Message: "Player statistic progression not found"}
)

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

// @summary Player Statistic Progression
// @description Get the player's statistic progression
// @router /api/v1/statistics/{statisticId}/players/{playerId} [GET]
// @produce json
// @param X-Game-ID header string true "Game ID responsible for the statistic"
// @param statisticId path string true "Statistic ID"
// @success 200 {object} PlayerStatisticProgression
// @failure 400,404,422,500 {object} ErrorResponse
func buildGetPlayerStatisticHandler(getPlayerProgressionFunc statistic.GetPlayerProgressionFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			statistic = c.Locals("statistic").(statistic.Statistic)
			playerID  = c.Params("playerId")
		)

		playerProgression, err := getPlayerProgressionFunc(c.Context(), statistic.ID, playerID)
		if err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(playerStatisticProgressionFromDomain(playerProgression))
	}
}
