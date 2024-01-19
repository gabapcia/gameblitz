package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gabarcia/game-blitz/internal/infra/logger/zap"
	"github.com/gabarcia/game-blitz/internal/quest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildStartPlayerQuestHandler(t *testing.T) {
	var (
		questID  = uuid.NewString()
		gameID   = uuid.NewString()
		playerID = uuid.NewString()

		expectedQuest = quest.Quest{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			ID:          questID,
			GameID:      gameID,
			Name:        "Test Quest",
			Description: "Start player quest handler unit test",
			Tasks: []quest.Task{
				{
					CreatedAt:             time.Now(),
					UpdatedAt:             time.Now(),
					ID:                    uuid.NewString(),
					Name:                  "Test Task",
					Description:           "Start player quest handler unit test",
					DependsOn:             make([]string, 0),
					RequiredForCompletion: true,
					Rule:                  `{"==": [{"var": "fields.bool"}, true]}`,
				},
			},
		}

		expectedPlayerProgression = quest.PlayerQuestProgression{
			StartedAt: time.Now(),
			UpdatedAt: time.Now(),
			PlayerID:  playerID,
			Quest:     expectedQuest,
			TasksProgression: []quest.PlayerTaskProgression{
				{
					StartedAt: time.Now(),
					UpdatedAt: time.Now(),
					Task:      expectedQuest.Tasks[0],
				},
			},
		}
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return expectedQuest, nil
			},
			StartQuestForPlayerFunc: func(ctx context.Context, q quest.Quest, playerID string) (quest.PlayerQuestProgression, error) {
				return expectedPlayerProgression, nil
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var body PlayerQuestProgression
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, expectedPlayerProgression.PlayerID, body.PlayerID)
		assert.Equal(t, expectedPlayerProgression.Quest.ID, body.Quest.ID)
		assert.Len(t, body.TasksProgression, 1)
	})

	t.Run("Quest Not Found", func(t *testing.T) {
		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return quest.Quest{}, quest.ErrQuestNotFound
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseQuestNotFound.Code, body.Code)
		assert.Equal(t, ErrorResponseQuestNotFound.Message, body.Message)
	})

	t.Run("Already Started", func(t *testing.T) {
		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return expectedQuest, nil
			},
			StartQuestForPlayerFunc: func(ctx context.Context, q quest.Quest, playerID string) (quest.PlayerQuestProgression, error) {
				return quest.PlayerQuestProgression{}, quest.ErrPlayerAlreadyStartedTheQuest
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponsePlayerAlreadyStartedTheQuest.Code, body.Code)
		assert.Equal(t, ErrorResponsePlayerAlreadyStartedTheQuest.Message, body.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return expectedQuest, nil
			},
			StartQuestForPlayerFunc: func(ctx context.Context, q quest.Quest, playerID string) (quest.PlayerQuestProgression, error) {
				return quest.PlayerQuestProgression{}, errors.New("any error")
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseInternalServerError.Code, body.Code)
		assert.Equal(t, ErrorResponseInternalServerError.Message, body.Message)
	})
}

func TestBuildGetPlayerQuestProgressionHandler(t *testing.T) {
	var (
		questID  = uuid.NewString()
		gameID   = uuid.NewString()
		playerID = uuid.NewString()

		expectedQuest = quest.Quest{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			ID:          questID,
			GameID:      gameID,
			Name:        "Test Quest",
			Description: "Start player quest handler unit test",
			Tasks: []quest.Task{
				{
					CreatedAt:             time.Now(),
					UpdatedAt:             time.Now(),
					ID:                    uuid.NewString(),
					Name:                  "Test Task",
					Description:           "Start player quest handler unit test",
					DependsOn:             make([]string, 0),
					RequiredForCompletion: true,
					Rule:                  `{"==": [{"var": "fields.bool"}, true]}`,
				},
			},
		}

		expectedPlayerProgression = quest.PlayerQuestProgression{
			StartedAt: time.Now(),
			UpdatedAt: time.Now(),
			PlayerID:  playerID,
			Quest:     expectedQuest,
			TasksProgression: []quest.PlayerTaskProgression{
				{
					StartedAt: time.Now(),
					UpdatedAt: time.Now(),
					Task:      expectedQuest.Tasks[0],
				},
			},
		}
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return expectedQuest, nil
			},
			GetPlayerQuestProgressionFunc: func(ctx context.Context, quest quest.Quest, playerID string) (quest.PlayerQuestProgression, error) {
				return expectedPlayerProgression, nil
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body PlayerQuestProgression
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, expectedPlayerProgression.PlayerID, body.PlayerID)
		assert.Equal(t, expectedPlayerProgression.Quest.ID, body.Quest.ID)
		assert.Len(t, body.TasksProgression, 1)
	})

	t.Run("Quest Not Found", func(t *testing.T) {
		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return quest.Quest{}, quest.ErrQuestNotFound
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseQuestNotFound.Code, body.Code)
		assert.Equal(t, ErrorResponseQuestNotFound.Message, body.Message)
	})

	t.Run("Not Started", func(t *testing.T) {
		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return expectedQuest, nil
			},
			GetPlayerQuestProgressionFunc: func(ctx context.Context, q quest.Quest, playerID string) (quest.PlayerQuestProgression, error) {
				return quest.PlayerQuestProgression{}, quest.ErrPlayerNotStartedTheQuest
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponsePlayerNotStartedTheQuest.Code, body.Code)
		assert.Equal(t, ErrorResponsePlayerNotStartedTheQuest.Message, body.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return expectedQuest, nil
			},
			GetPlayerQuestProgressionFunc: func(ctx context.Context, q quest.Quest, playerID string) (quest.PlayerQuestProgression, error) {
				return quest.PlayerQuestProgression{}, errors.New("any error")
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseInternalServerError.Code, body.Code)
		assert.Equal(t, ErrorResponseInternalServerError.Message, body.Message)
	})
}

func TestBuildUpdatePlayerQuestProgressionHandler(t *testing.T) {
	var (
		questID  = uuid.NewString()
		gameID   = uuid.NewString()
		playerID = uuid.NewString()

		expectedQuest = quest.Quest{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			ID:          questID,
			GameID:      gameID,
			Name:        "Test Quest",
			Description: "Start player quest handler unit test",
			Tasks: []quest.Task{
				{
					CreatedAt:             time.Now(),
					UpdatedAt:             time.Now(),
					ID:                    uuid.NewString(),
					Name:                  "Test Task",
					Description:           "Start player quest handler unit test",
					DependsOn:             make([]string, 0),
					RequiredForCompletion: true,
					Rule:                  `{"==": [{"var": "fields.bool"}, true]}`,
				},
			},
		}

		expectedPlayerProgression = quest.PlayerQuestProgression{
			StartedAt: time.Now(),
			UpdatedAt: time.Now(),
			PlayerID:  playerID,
			Quest:     expectedQuest,
			TasksProgression: []quest.PlayerTaskProgression{
				{
					StartedAt: time.Now(),
					UpdatedAt: time.Now(),
					Task:      expectedQuest.Tasks[0],
				},
			},
		}
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return expectedQuest, nil
			},
			UpdatePlayerQuestProgressionFunc: func(ctx context.Context, q quest.Quest, playerID, taskDataToCheck string) (quest.PlayerQuestProgression, error) {
				data := expectedPlayerProgression
				data.CompletedAt = time.Now()
				data.TasksProgression[0].CompletedAt = time.Now()
				return data, nil
			},
		})

		data, err := json.Marshal(map[string]string{
			"data": `{"fields": {"bool": true}}`,
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), bytes.NewBuffer(data))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body PlayerQuestProgression
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, expectedPlayerProgression.PlayerID, body.PlayerID)
		assert.Equal(t, expectedPlayerProgression.Quest.ID, body.Quest.ID)
		assert.NotEmpty(t, body.CompletedAt)
		if assert.Len(t, body.TasksProgression, 1) {
			assert.NotEmpty(t, body.TasksProgression[0].CompletedAt)
		}
	})

	t.Run("Quest Not Found", func(t *testing.T) {
		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return quest.Quest{}, quest.ErrQuestNotFound
			},
		})

		data, err := json.Marshal(map[string]string{
			"data": `{"fields": {"bool": true}}`,
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), bytes.NewBuffer(data))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseQuestNotFound.Code, body.Code)
		assert.Equal(t, ErrorResponseQuestNotFound.Message, body.Message)
	})

	t.Run("Already Completed", func(t *testing.T) {
		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return expectedQuest, nil
			},
			UpdatePlayerQuestProgressionFunc: func(ctx context.Context, q quest.Quest, playerID, taskDataToCheck string) (quest.PlayerQuestProgression, error) {
				return quest.PlayerQuestProgression{}, quest.ErrPlayerQuestAlreadyCompleted
			},
		})

		data, err := json.Marshal(map[string]string{
			"data": `{"fields": {"bool": true}}`,
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), bytes.NewBuffer(data))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponsePlayerQuestAlreadyFinished.Code, body.Code)
		assert.Equal(t, ErrorResponsePlayerQuestAlreadyFinished.Message, body.Message)
	})

	t.Run("Not Started", func(t *testing.T) {
		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return expectedQuest, nil
			},
			UpdatePlayerQuestProgressionFunc: func(ctx context.Context, q quest.Quest, playerID, taskDataToCheck string) (quest.PlayerQuestProgression, error) {
				return quest.PlayerQuestProgression{}, quest.ErrPlayerNotStartedTheQuest
			},
		})

		data, err := json.Marshal(map[string]string{
			"data": `{"fields": {"bool": true}}`,
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), bytes.NewBuffer(data))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponsePlayerNotStartedTheQuest.Code, body.Code)
		assert.Equal(t, ErrorResponsePlayerNotStartedTheQuest.Message, body.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		app := App(Config{
			GetQuestByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (quest.Quest, error) {
				return expectedQuest, nil
			},
			UpdatePlayerQuestProgressionFunc: func(ctx context.Context, q quest.Quest, playerID, taskDataToCheck string) (quest.PlayerQuestProgression, error) {
				return quest.PlayerQuestProgression{}, errors.New("any error")
			},
		})

		data, err := json.Marshal(map[string]string{
			"data": `{"fields": {"bool": true}}`,
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/quests/%s/players/%s", questID, playerID), bytes.NewBuffer(data))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseInternalServerError.Code, body.Code)
		assert.Equal(t, ErrorResponseInternalServerError.Message, body.Message)
	})
}
