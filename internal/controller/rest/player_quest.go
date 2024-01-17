package rest

import (
	"net/http"
	"time"

	"github.com/gabarcia/metagaming-api/internal/quest"

	"github.com/gofiber/fiber/v2"
)

type (
	PlayerQuestTaskProgression struct {
		StartedAt   time.Time  `json:"startedAt"`             // Time the player started the task
		UpdatedAt   time.Time  `json:"updatedAt"`             // Last time the player updated the task progression
		Task        Task       `json:"task"`                  // Task config data
		CompletedAt *time.Time `json:"completedAt,omitempty"` // Time the player completed the task
	}

	PlayerQuestProgression struct {
		StartedAt        time.Time                    `json:"startedAt"`             // Time the player started the quest
		UpdatedAt        time.Time                    `json:"updatedAt"`             // Last time the player updated the quest progression
		PlayerID         string                       `json:"playerId"`              // Player's ID
		Quest            Quest                        `json:"quest"`                 // Quest Config Data
		CompletedAt      *time.Time                   `json:"completedAt,omitempty"` // Time the player completed the quest
		TasksProgression []PlayerQuestTaskProgression `json:"tasksProgression"`      // Tasks progression
	}
)

func playerQuestProgressionFromDomain(p quest.PlayerQuestProgression) PlayerQuestProgression {
	tasksProgression := make([]PlayerQuestTaskProgression, len(p.TasksProgression))
	for i, tp := range p.TasksProgression {
		var completedAt *time.Time
		if !tp.CompletedAt.IsZero() {
			tmp := tp.CompletedAt
			completedAt = &tmp
		}

		tasksProgression[i] = PlayerQuestTaskProgression{
			StartedAt:   tp.StartedAt,
			UpdatedAt:   tp.UpdatedAt,
			Task:        taskFromDomain(tp.Task),
			CompletedAt: completedAt,
		}
	}

	var completedAt *time.Time
	if !p.CompletedAt.IsZero() {
		completedAt = &p.CompletedAt
	}

	return PlayerQuestProgression{
		StartedAt:        p.StartedAt,
		UpdatedAt:        p.UpdatedAt,
		PlayerID:         p.PlayerID,
		Quest:            questFromDomain(p.Quest),
		CompletedAt:      completedAt,
		TasksProgression: tasksProgression,
	}
}

// @summary Start Player Quest Progression
// @description Start a player's quest progression
// @router /api/v1/quests/{questId}/players/{playerId} [POST]
// @produce json
// @param X-Game-ID header string true "Game ID responsible for the quest"
// @param questId path string true "Quest ID"
// @param playerId path string true "Player ID"
// @success 201 {object} PlayerQuestProgression
// @failure 404,422,500 {object} ErrorResponse
func buildStartPlayerQuestHandler(startQuestForPlayerFunc quest.StartQuestForPlayerFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			quest    = c.Locals("quest").(quest.Quest)
			playerID = c.Params("playerId")
		)

		progression, err := startQuestForPlayerFunc(c.Context(), quest, playerID)
		if err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(playerQuestProgressionFromDomain(progression))
	}
}

// @summary Get Player Quest Progression
// @description Get a player's quest progression
// @router /api/v1/quests/{questId}/players/{playerId} [GET]
// @produce json
// @param X-Game-ID header string true "Game ID responsible for the quest"
// @param questId path string true "Quest ID"
// @param playerId path string true "Player ID"
// @success 200 {object} PlayerQuestProgression
// @failure 404,422,500 {object} ErrorResponse
func buildGetPlayerQuestProgressionHandler(getPlayerQuestProgressionFunc quest.GetPlayerQuestProgressionFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			quest    = c.Locals("quest").(quest.Quest)
			playerID = c.Params("playerId")
		)

		progression, err := getPlayerQuestProgressionFunc(c.Context(), quest, playerID)
		if err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(playerQuestProgressionFromDomain(progression))
	}
}