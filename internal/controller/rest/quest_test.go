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

	"github.com/gabarcia/metagaming-api/internal/infra/logger/zap"
	"github.com/gabarcia/metagaming-api/internal/quest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildBuildCreateQuestHanlder(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		gameID := uuid.NewString()
		app := App(Config{
			CreateQuestFunc: func(ctx context.Context, data quest.NewQuestData) (quest.Quest, error) {
				tasks := make([]quest.Task, len(data.Tasks))
				for i := range data.Tasks {
					tasks[i] = quest.Task{ID: uuid.NewString()}
				}

				return quest.Quest{
					ID:     uuid.NewString(),
					GameID: data.GameID,
					Tasks:  tasks,
				}, nil
			},
		})

		body, err := json.Marshal(map[string]any{
			"name":        "Test Create Quest",
			"description": "Test create quest handler unit test",
			"tasks": []map[string]any{
				{
					"name":        "Test Task #0",
					"description": "Test task description",
					"rule":        `{"==": [{"var": {"fields.bool"}}, true]}`,
				},
				{
					"name":        "Test Task #0",
					"description": "Test task description",
					"dependsOn":   0,
					"rule":        `{"==": [{"var": {"fields.bool"}}, true]}`,
				},
			},
			"tasksValidators": []string{
				`{"fields": {"bool": true}}`,
				`{"fields": {"bool": true}}`,
			},
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/quests", bytes.NewBuffer(body))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var data Quest
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.NotEmpty(t, data.ID)
		assert.Equal(t, gameID, data.GameID)
		assert.NotEmpty(t, data.Tasks)
	})

	t.Run("Validation Error", func(t *testing.T) {
		gameID := uuid.NewString()
		app := App(Config{
			CreateQuestFunc: quest.BuildCreateQuestFunc(nil),
		})

		body, err := json.Marshal(map[string]any{
			"tasks": []map[string]any{
				{
					"rule": `{`,
				},
			},
			"tasksValidators": []string{},
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/quests", bytes.NewBuffer(body))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseQuestInvalid.Code, data.Code)
		assert.Equal(t, ErrorResponseQuestInvalid.Message, data.Message)
	})

	t.Run("Missing Game ID", func(t *testing.T) {
		app := App(Config{})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/quests", bytes.NewBufferString("{}"))

		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseQuestInvalidGameID.Code, data.Code)
		assert.Equal(t, ErrorResponseQuestInvalidGameID.Message, data.Message)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		var (
			gameID = uuid.NewString()
			app    = App(Config{})
			req    = httptest.NewRequest(http.MethodPost, "/api/v1/quests", bytes.NewBufferString("{"))
		)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseInvalidRequestBody.Code, data.Code)
		assert.Equal(t, ErrorResponseInvalidRequestBody.Message, data.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		gameID := uuid.NewString()
		app := App(Config{
			CreateQuestFunc: func(ctx context.Context, data quest.NewQuestData) (quest.Quest, error) {
				return quest.Quest{}, errors.New("any error")
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/quests", bytes.NewBufferString("{}"))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseInternalServerError.Code, data.Code)
		assert.Equal(t, ErrorResponseInternalServerError.Message, data.Message)
	})
}

func TestBuildBuildDeleteQuestHanlder(t *testing.T) {
	var (
		questID = uuid.NewString()
		gameID  = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			SoftDeleteQuestFunc: func(ctx context.Context, questID, gameID string) error {
				return nil
			},
		})

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/quests/%s", questID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Invalid Quest ID", func(t *testing.T) {
		app := App(Config{
			SoftDeleteQuestFunc: func(ctx context.Context, questID, gameID string) error {
				return quest.ErrInvalidQuestID
			},
		})

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/quests/invalid-id", nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseQuestInvalidID.Code, data.Code)
		assert.Equal(t, ErrorResponseQuestInvalidID.Message, data.Message)
	})

	t.Run("Missing Game ID", func(t *testing.T) {
		app := App(Config{})

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/quests/%s", questID), nil)

		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseQuestInvalidGameID.Code, data.Code)
		assert.Equal(t, ErrorResponseQuestInvalidGameID.Message, data.Message)
	})

	t.Run("Not Found", func(t *testing.T) {
		app := App(Config{
			SoftDeleteQuestFunc: func(ctx context.Context, questID, gameID string) error {
				return quest.ErrQuestNotFound
			},
		})

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/quests/%s", questID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseQuestNotFound.Code, data.Code)
		assert.Equal(t, ErrorResponseQuestNotFound.Message, data.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		app := App(Config{
			SoftDeleteQuestFunc: func(ctx context.Context, questID, gameID string) error {
				return errors.New("any error")
			},
		})

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/quests/%s", questID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseInternalServerError.Code, data.Code)
		assert.Equal(t, ErrorResponseInternalServerError.Message, data.Message)
	})
}
