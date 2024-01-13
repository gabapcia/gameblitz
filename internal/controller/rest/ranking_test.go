package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
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

	t.Run("Invalid Request Body", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/leaderboards/%s/ranking/%s", leaderboardID, playerID), bytes.NewBufferString(`{`))

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

	t.Run("Missing Game ID", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
			},
		})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/leaderboards/%s/ranking/%s", leaderboardID, playerID), bytes.NewBufferString(`{"value": 100.0}`))

		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseLeaderboardInvalidGameID.Code, body.Code)
		assert.Equal(t, ErrorResponseLeaderboardInvalidGameID.Message, body.Message)
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

func TestBuildGetRankingHandler(t *testing.T) {
	var (
		leaderboardID = uuid.NewString()
		gameID        = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
			},
			RankingFunc: func(ctx context.Context, leaderboard leaderboard.Leaderboard, page, limit int64) ([]ranking.Rank, error) {
				rankings := make([]ranking.Rank, 0)
				for i := 0; i < 10; i++ {
					rankings = append(rankings, ranking.Rank{
						LeaderboardID: uuid.NewString(),
						PlayerID:      uuid.NewString(),
						Position:      int64(i),
						Value:         rand.Float64(),
					})
				}

				return rankings, nil
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/leaderboards/%s/ranking", leaderboardID), nil)

		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var data []Rank
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)
	})

	t.Run("Missing Game ID", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/leaderboards/%s/ranking", leaderboardID), nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseLeaderboardInvalidGameID.Code, body.Code)
		assert.Equal(t, ErrorResponseLeaderboardInvalidGameID.Message, body.Message)
	})

	t.Run("Leaderboard Not Found", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{}, leaderboard.ErrLeaderboardNotFound
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/leaderboards/%s/ranking", leaderboardID), nil)

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

	t.Run("Invalid Page Number", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
			},
			RankingFunc: func(ctx context.Context, leaderboard leaderboard.Leaderboard, page, limit int64) ([]ranking.Rank, error) {
				return nil, ranking.ErrInvalidPageNumber
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/leaderboards/%s/ranking", leaderboardID), nil)

		req.Header.Set(gameIDHeader, gameID)

		q := req.URL.Query()
		q.Set("page", "-1")
		req.URL.RawQuery = q.Encode()

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseRankingPageNumber.Code, body.Code)
		assert.Equal(t, ErrorResponseRankingPageNumber.Message, body.Message)
	})

	t.Run("Invalid Limit Number", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
			},
			RankingFunc: func(ctx context.Context, leaderboard leaderboard.Leaderboard, page, limit int64) ([]ranking.Rank, error) {
				return nil, ranking.ErrInvalidLimitNumber
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/leaderboards/%s/ranking", leaderboardID), nil)

		req.Header.Set(gameIDHeader, gameID)

		q := req.URL.Query()
		q.Set("limit", "-1")
		req.URL.RawQuery = q.Encode()

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseRankingLimitNumber.Code, body.Code)
		assert.Equal(t, ErrorResponseRankingLimitNumber.Message, body.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
			},
			RankingFunc: func(ctx context.Context, leaderboard leaderboard.Leaderboard, page, limit int64) ([]ranking.Rank, error) {
				return nil, errors.New("any error")
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/leaderboards/%s/ranking", leaderboardID), nil)

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
