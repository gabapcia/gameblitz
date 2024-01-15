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
	"github.com/gabarcia/metagaming-api/internal/leaderboard"
	"github.com/gabarcia/metagaming-api/internal/statistic"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTuildUpsertPlayerStatisticHandler(t *testing.T) {
	var (
		statisticID = uuid.NewString()
		gameID      = uuid.NewString()
		playerID    = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{ID: id, GameID: gameID}, nil
			},
			UpsertPlayerStatisticProgressionFunc: func(ctx context.Context, statistic statistic.Statistic, playerID string, value float64) error {
				return nil
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), bytes.NewBufferString(`{"value": 100.0}`))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		app := App(Config{
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{ID: id, GameID: gameID}, nil
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), bytes.NewBufferString(`{`))

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

	t.Run("Statistic Not Found", func(t *testing.T) {
		app := App(Config{
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{}, statistic.ErrStatisticNotFound
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), bytes.NewBufferString(`{"value": 100.0}`))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseStatisticNotFound.Code, body.Code)
		assert.Equal(t, ErrorResponseStatisticNotFound.Message, body.Message)
	})

	t.Run("Missing Game ID", func(t *testing.T) {
		app := App(Config{
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{ID: id, GameID: gameID}, nil
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), bytes.NewBufferString(`{"value": 100.0}`))

		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseStatisticInvalidGameID.Code, body.Code)
		assert.Equal(t, ErrorResponseStatisticInvalidGameID.Message, body.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		app := App(Config{
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{ID: id, GameID: gameID}, nil
			},
			UpsertPlayerRankFunc: func(ctx context.Context, leaderboard leaderboard.Leaderboard, playerID string, value float64) error {
				return errors.New("any error")
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), bytes.NewBufferString(`{"value": 100.0}`))

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
