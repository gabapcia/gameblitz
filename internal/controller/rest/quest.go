package rest

import (
	"net/http"
	"time"

	"github.com/gabarcia/metagaming-api/internal/quest"
	"github.com/gofiber/fiber/v2"
)

type CreateQuestReq struct {
	Name        string `json:"name"`        // Quest name
	Description string `json:"description"` // Quest details
	Tasks       []struct {
		Name        string `json:"name"`        // Task name
		Description string `json:"description"` // Task details
		DependsOn   *int   `json:"dependsOn"`   // Array index of the task that needs to be completed before this one can be started
		Rule        string `json:"rule"`        // Task completion logic as JsonLogic. See https://jsonlogic.com/
	} `json:"tasks"` // Quest task list
	TasksValidators []string `json:"tasksValidators"` // Quest task list success validation data
}

type Quest struct {
	CreatedAt   time.Time `json:"createdAt"`   // Time that the quest was created
	UpdatedAt   time.Time `json:"updatedAt"`   // Last time that the quest was updated
	ID          string    `json:"id"`          // Quest ID
	GameID      string    `json:"gameId"`      // ID of the game responsible for the quest
	Name        string    `json:"name"`        // Quest name
	Description string    `json:"description"` // Quest details
	Tasks       []Task    `json:"tasks"`       // Quest task list
}

func (q CreateQuestReq) toDomain(gameID string) quest.NewQuestData {
	tasks := make([]quest.NewTaskData, len(q.Tasks))
	for i, t := range q.Tasks {
		tasks[i] = quest.NewTaskData{
			Name:        t.Name,
			Description: t.Description,
			DependsOn:   t.DependsOn,
			Rule:        t.Rule,
		}
	}

	return quest.NewQuestData{
		GameID:          gameID,
		Name:            q.Name,
		Description:     q.Description,
		Tasks:           tasks,
		TasksValidators: q.TasksValidators,
	}
}

func questFromDomain(q quest.Quest) Quest {
	tasks := make([]Task, len(q.Tasks))
	for i, task := range q.Tasks {
		tasks[i] = taskFromDomain(task)
	}

	return Quest{
		CreatedAt:   q.CreatedAt,
		UpdatedAt:   q.UpdatedAt,
		ID:          q.ID,
		GameID:      q.GameID,
		Name:        q.Name,
		Description: q.Description,
		Tasks:       tasks,
	}
}

var (
	ErrorResponseQuestInvalidGameID = ErrorResponse{Code: "3.0", Message: "Invalid Game ID"}
	ErrorResponseQuestInvalid       = ErrorResponse{Code: "3.1", Message: "Invalid quest data"}
	ErrorResponseQuestNotFound      = ErrorResponse{Code: "3.2", Message: "Quest not found"}
	ErrorResponseQuestInvalidID     = ErrorResponse{Code: "3.3", Message: "Invalid quest id"}
)

// @summary Create Quest
// @description Create a quest and its tasks
// @router /api/v1/quests [POST]
// @accept json
// @produce json
// @param X-Game-ID header string true "Game ID responsible for the leaderboard"
// @param NewQuestData body CreateQuestReq true "New quest config data"
// @success 201 {object} Quest
// @failure 400,422,500 {object} ErrorResponse
func buildBuildCreateQuestHanlder(createQuestFunc quest.CreateQuestFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		gameID := string(c.Request().Header.Peek(gameIDHeader))
		if gameID == "" {
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseQuestInvalidGameID)
		}

		var body CreateQuestReq
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		quest, err := createQuestFunc(c.Context(), body.toDomain(gameID))
		if err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(questFromDomain(quest))
	}
}

// @summary Get Quest By ID
// @description Get a quest and its tasks
// @router /api/v1/quests/{questId} [GET]
// @produce json
// @param X-Game-ID header string true "Game ID responsible for the leaderboard"
// @param questId path string true "Quest ID"
// @success 200 {object} Quest
// @failure 404,422,500 {object} ErrorResponse
func buildBuildGetQuestHanlder(getQuestByIDAndGameID quest.GetQuestByIDAndGameIDFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			questID = c.Params("questId")
			gameID  = string(c.Request().Header.Peek(gameIDHeader))
		)

		if gameID == "" {
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseQuestInvalidGameID)
		}

		quest, err := getQuestByIDAndGameID(c.Context(), questID, gameID)
		if err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(questFromDomain(quest))
	}
}

// @summary Delete Quest
// @description Delete a quest and its tasks
// @router /api/v1/quests/{questId} [DELETE]
// @param X-Game-ID header string true "Game ID responsible for the leaderboard"
// @param questId path string true "Quest ID"
// @success 204
// @failure 404,422,500 {object} ErrorResponse
func buildBuildDeleteQuestHanlder(softDeleteQuestFunc quest.SoftDeleteQuestFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			questID = c.Params("questId")
			gameID  = string(c.Request().Header.Peek(gameIDHeader))
		)

		if gameID == "" {
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseQuestInvalidGameID)
		}

		if err := softDeleteQuestFunc(c.Context(), questID, gameID); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	}
}
