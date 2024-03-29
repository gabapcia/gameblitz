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

	"github.com/gabapcia/gameblitz/internal/auth"
	"github.com/gabapcia/gameblitz/internal/infra/logger/zap"
	"github.com/gabapcia/gameblitz/internal/leaderboard"
	"github.com/gabapcia/gameblitz/internal/statistic"
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
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{ID: id, GameID: gameID}, nil
			},
			UpsertPlayerStatisticProgressionFunc: func(ctx context.Context, statistic statistic.Statistic, playerID string, value float64) error {
				return nil
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), bytes.NewBufferString(`{"value": 100.0}`))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{ID: id, GameID: gameID}, nil
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), bytes.NewBufferString(`{`))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

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
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{}, statistic.ErrStatisticNotFound
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), bytes.NewBufferString(`{"value": 100.0}`))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseStatisticNotFound.Code, body.Code)
		assert.Equal(t, ErrorResponseStatisticNotFound.Message, body.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{ID: id, GameID: gameID}, nil
			},
			UpsertPlayerRankFunc: func(ctx context.Context, leaderboard leaderboard.Leaderboard, playerID string, value float64) error {
				return errors.New("any error")
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), bytes.NewBufferString(`{"value": 100.0}`))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

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

func TestBuildGetPlayerStatisticHandler(t *testing.T) {
	var (
		statisticID = uuid.NewString()
		gameID      = uuid.NewString()
		playerID    = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{ID: id, GameID: gameID}, nil
			},
			GetPlayerStatisticProgressionFunc: func(ctx context.Context, statisticID, playerID string) (statistic.PlayerProgression, error) {
				return statistic.PlayerProgression{PlayerID: playerID, StatisticID: statisticID}, nil
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), nil)

		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var data PlayerStatisticProgression
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, statisticID, data.StatisticID)
		assert.Equal(t, playerID, data.PlayerID)
	})

	t.Run("Statistic Not Found", func(t *testing.T) {
		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{}, statistic.ErrStatisticNotFound
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), nil)

		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseStatisticNotFound.Code, body.Code)
		assert.Equal(t, ErrorResponseStatisticNotFound.Message, body.Message)
	})

	t.Run("Player Statistic Progression Not Found", func(t *testing.T) {
		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{ID: id, GameID: gameID}, nil
			},
			GetPlayerStatisticProgressionFunc: func(ctx context.Context, statisticID, playerID string) (statistic.PlayerProgression, error) {
				return statistic.PlayerProgression{}, statistic.ErrPlayerStatisticNotFound
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), nil)

		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponsePlayerStatisticNotFound.Code, body.Code)
		assert.Equal(t, ErrorResponsePlayerStatisticNotFound.Message, body.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{ID: id, GameID: gameID}, nil
			},
			GetPlayerStatisticProgressionFunc: func(ctx context.Context, statisticID, playerID string) (statistic.PlayerProgression, error) {
				return statistic.PlayerProgression{}, errors.New("any error")
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/statistics/%s/players/%s", statisticID, playerID), nil)

		req.Header.Set("Authorization", uuid.NewString())

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
