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
	"github.com/gabarcia/metagaming-api/internal/ranking"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildUpsertPlayerRankHandler(t *testing.T) {
	var (
		leaderboardID = uuid.NewString()
		gameID        = uuid.NewString()
		playerID      = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
			},
			UpsertPlayerRankFunc: func(ctx context.Context, leaderboard leaderboard.Leaderboard, playerID string, value float64) error {
				return nil
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/leaderboards/%s/ranking/%s", leaderboardID, playerID), bytes.NewBufferString(`{"value": 100.0}`))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Leaderboard Not Found", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{}, leaderboard.ErrLeaderboardNotFound
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/leaderboards/%s/ranking/%s", leaderboardID, playerID), bytes.NewBufferString(`{"value": 100.0}`))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseLeaderboardNotFound.Code, body.Code)
		assert.Equal(t, ErrorResponseLeaderboardNotFound.Message, body.Message)
	})

	t.Run("Leaderboard Closed", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
			},
			UpsertPlayerRankFunc: func(ctx context.Context, leaderboard leaderboard.Leaderboard, playerID string, value float64) error {
				return ranking.ErrLeaderboardClosed
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/leaderboards/%s/ranking/%s", leaderboardID, playerID), bytes.NewBufferString(`{"value": 100.0}`))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseLeaderboardClosed.Code, body.Code)
		assert.Equal(t, ErrorResponseLeaderboardClosed.Message, body.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
			},
			UpsertPlayerRankFunc: func(ctx context.Context, leaderboard leaderboard.Leaderboard, playerID string, value float64) error {
				return errors.New("any error")
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/leaderboards/%s/ranking/%s", leaderboardID, playerID), bytes.NewBufferString(`{"value": 100.0}`))

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
