package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gabapcia/gameblitz/internal/auth"
	"github.com/gabapcia/gameblitz/internal/infra/logger/zap"
	"github.com/gabapcia/gameblitz/internal/quest"
	"github.com/gofiber/fiber/v2"
)

type CreateQuestReq struct {
	Name        string `json:"name"`        // Quest name
	Description string `json:"description"` // Quest details
	Tasks       []struct {
		Name                  string `json:"name"`                  // Task name
		Description           string `json:"description"`           // Task details
		DependsOn             []int  `json:"dependsOn"`             // List of array indexes of the tasks that needs to be completed before this one can be started
		RequiredForCompletion *bool  `json:"requiredForCompletion"` // Is this task required for the quest completion? Defaults to `true`
		Rule                  string `json:"rule"`                  // Task completion logic as JsonLogic. See https://jsonlogic.com/
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
		requiredForCompletion := true
		if t.RequiredForCompletion != nil {
			requiredForCompletion = *t.RequiredForCompletion
		}

		tasks[i] = quest.NewTaskData{
			Name:                  t.Name,
			Description:           t.Description,
			DependsOn:             t.DependsOn,
			RequiredForCompletion: requiredForCompletion,
			Rule:                  t.Rule,
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
	ErrorResponseQuestInvalid   = ErrorResponse{Code: "3.0", Message: "Invalid quest data"}
	ErrorResponseQuestNotFound  = ErrorResponse{Code: "3.1", Message: "Quest not found"}
	ErrorResponseQuestInvalidID = ErrorResponse{Code: "3.2", Message: "Invalid quest id"}
)

func buildGetQuestMiddleware(cache fiber.Storage, expiration time.Duration, getQuestByIDAndGameIDFunc quest.GetQuestByIDAndGameIDFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			id       = c.Params("questId")
			claims   = c.Locals("claims").(auth.Claims)
			cacheKey = fmt.Sprintf("GetQuestMiddleware:%s:%s", id, claims.GameID)
		)

		if cache != nil {
			data, err := cache.Get(cacheKey)
			if err != nil {
				zap.Error(err, "get cache error")
			} else if data != nil {
				var quest quest.Quest
				if err = json.Unmarshal(data, &quest); err != nil {
					zap.Error(err, "unmarshal cached quest error")
				} else {
					c.Locals("quest", quest)
					return c.Next()
				}
			}
		}

		quest, err := getQuestByIDAndGameIDFunc(c.Context(), id, claims.GameID)
		if err != nil {
			return err
		}

		if cache != nil {
			data, err := json.Marshal(quest)
			if err != nil {
				zap.Error(err, "marshal quest cache error")
			} else {
				if err = cache.Set(cacheKey, data, expiration); err != nil {
					zap.Error(err, "unable to cache quest")
				}
			}
		}

		c.Locals("quest", quest)
		return c.Next()
	}
}

// @summary Create Quest
// @description Create a quest and its tasks
// @router /api/v1/quests [POST]
// @accept json
// @produce json
// @param Authorization header string true "Game's JWT authorization"
// @param NewQuestData body CreateQuestReq true "New quest config data"
// @success 201 {object} Quest
// @failure 400,422,500 {object} ErrorResponse
func buildCreateQuestHanlder(createQuestFunc quest.CreateQuestFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(auth.Claims)

		var body CreateQuestReq
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		quest, err := createQuestFunc(c.Context(), body.toDomain(claims.GameID))
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
// @param Authorization header string true "Game's JWT authorization"
// @param questId path string true "Quest ID"
// @success 200 {object} Quest
// @failure 404,422,500 {object} ErrorResponse
func buildGetQuestHanlder(getQuestByIDAndGameID quest.GetQuestByIDAndGameIDFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			questID = c.Params("questId")
			claims  = c.Locals("claims").(auth.Claims)
		)

		quest, err := getQuestByIDAndGameID(c.Context(), questID, claims.GameID)
		if err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(questFromDomain(quest))
	}
}

// @summary Delete Quest
// @description Delete a quest and its tasks
// @router /api/v1/quests/{questId} [DELETE]
// @param Authorization header string true "Game's JWT authorization"
// @param questId path string true "Quest ID"
// @success 204
// @failure 404,422,500 {object} ErrorResponse
func buildDeleteQuestHanlder(softDeleteQuestFunc quest.SoftDeleteQuestFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			questID = c.Params("questId")
			claims  = c.Locals("claims").(auth.Claims)
		)

		if err := softDeleteQuestFunc(c.Context(), questID, claims.GameID); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	}
}
