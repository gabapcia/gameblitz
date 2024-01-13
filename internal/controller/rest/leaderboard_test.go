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

	"github.com/gabarcia/metagaming-api/internal/infra/logger/zap"
	"github.com/gabarcia/metagaming-api/internal/leaderboard"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildGetLeaderboardMiddleware(t *testing.T) {
	var (
		leaderboardID = uuid.NewString()
		gameID        = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		getLeaderboardMiddleware := buildGetLeaderboardMiddleware(nil, time.Minute, func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
			return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
		})

		app := fiber.New()
		app.Get("/:leaderboardId", getLeaderboardMiddleware, func(c *fiber.Ctx) error {
			leaderboard := c.Locals("leaderboard")
			return c.Status(http.StatusOK).JSON(leaderboard)
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", leaderboardID), nil)

		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body Leaderboard
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, leaderboardID, body.ID)
		assert.Equal(t, gameID, body.GameID)
	})

	t.Run("Missing Game ID", func(t *testing.T) {
		getLeaderboardMiddleware := buildGetLeaderboardMiddleware(nil, time.Minute, func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
			return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
		})

		app := fiber.New()
		app.Get("/:leaderboardId", getLeaderboardMiddleware, func(c *fiber.Ctx) error {
			leaderboard := c.Locals("leaderboard")
			return c.Status(http.StatusOK).JSON(leaderboard)
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", leaderboardID), nil)

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
		getLeaderboardMiddleware := buildGetLeaderboardMiddleware(nil, time.Minute, func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
			return leaderboard.Leaderboard{}, leaderboard.ErrLeaderboardNotFound
		})

		app := fiber.New(fiber.Config{ErrorHandler: buildErrorHandler()})
		app.Get("/:leaderboardId", getLeaderboardMiddleware, func(c *fiber.Ctx) error {
			leaderboard := c.Locals("leaderboard")
			return c.Status(http.StatusOK).JSON(leaderboard)
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", leaderboardID), nil)

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
}

func TestBuildCreateLeaderboardHandler(t *testing.T) {
	var (
		gameID          = uuid.NewString()
		name            = "Test Leaderboard"
		description     = "Test create leaderboard request"
		startAt         = time.Now().Format(time.RFC3339)
		endAt           = time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		aggregationMode = "MAX"
		ordering        = "DESC"
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			CreateLeaderboardFunc: leaderboard.BuildCreateFunc(func(ctx context.Context, data leaderboard.NewLeaderboardData) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{
					ID:              uuid.NewString(),
					GameID:          data.GameID,
					Name:            data.Name,
					Description:     data.Description,
					StartAt:         data.StartAt,
					EndAt:           data.EndAt,
					AggregationMode: data.AggregationMode,
					Ordering:        data.Ordering,
				}, nil
			}),
		})

		reqBody, err := json.Marshal(map[string]any{
			"name":            name,
			"description":     description,
			"startAt":         startAt,
			"endAt":           endAt,
			"aggregationMode": aggregationMode,
			"ordering":        ordering,
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/leaderboards", bytes.NewBuffer(reqBody))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var data Leaderboard
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.NotEmpty(t, data.ID)
		assert.Equal(t, gameID, data.GameID)
		assert.Equal(t, name, data.Name)
		assert.Equal(t, description, data.Description)
		assert.Equal(t, startAt, data.StartAt.Format(time.RFC3339))
		assert.Equal(t, endAt, data.EndAt.Format(time.RFC3339))
		assert.Equal(t, aggregationMode, data.AggregationMode)
		assert.Equal(t, ordering, data.Ordering)
	})

	t.Run("Validation Error", func(t *testing.T) {
		app := App(Config{
			CreateLeaderboardFunc: leaderboard.BuildCreateFunc(func(ctx context.Context, data leaderboard.NewLeaderboardData) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{}, nil
			}),
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/leaderboards", bytes.NewBufferString(`{}`))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseLeaderboardInvalid.Code, data.Code)
		assert.Equal(t, ErrorResponseLeaderboardInvalid.Message, data.Message)
		assert.NotEmpty(t, data.Details)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		app := App(Config{
			CreateLeaderboardFunc: leaderboard.BuildCreateFunc(func(ctx context.Context, data leaderboard.NewLeaderboardData) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{ID: uuid.NewString()}, nil
			}),
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/leaderboards", bytes.NewBufferString(`{`))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

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

		app := App(Config{
			CreateLeaderboardFunc: leaderboard.BuildCreateFunc(func(ctx context.Context, data leaderboard.NewLeaderboardData) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{}, errors.New("any error")
			}),
		})

		reqBody, err := json.Marshal(map[string]any{
			"name":            name,
			"description":     description,
			"startAt":         startAt,
			"endAt":           endAt,
			"aggregationMode": aggregationMode,
			"ordering":        ordering,
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/leaderboards", bytes.NewBuffer(reqBody))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(gameIDHeader, gameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseInternalServerError.Code, data.Code)
		assert.Equal(t, ErrorResponseInternalServerError.Message, data.Message)
		assert.Empty(t, data.Details)
	})
}

func TestBuildGetLeaderboardHandler(t *testing.T) {
	var (
		expectedID     = uuid.NewString()
		expectedGameID = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: leaderboard.BuildGetByIDAndGameIDFunc(func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{ID: id, GameID: gameID}, nil
			}),
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/leaderboards/%s", expectedID), nil)

		req.Header.Set(gameIDHeader, expectedGameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var data Leaderboard
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, expectedID, data.ID)
		assert.Equal(t, expectedGameID, data.GameID)
	})

	t.Run("Missing Game ID", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: leaderboard.BuildGetByIDAndGameIDFunc(func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{}, nil
			}),
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/leaderboards/%s", expectedID), nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseLeaderboardInvalidGameID.Code, data.Code)
		assert.Equal(t, ErrorResponseLeaderboardInvalidGameID.Message, data.Message)
	})

	t.Run("Not Found", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: leaderboard.BuildGetByIDAndGameIDFunc(func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{}, leaderboard.ErrLeaderboardNotFound
			}),
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/leaderboards/%s", expectedID), nil)

		req.Header.Set(gameIDHeader, expectedGameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseLeaderboardNotFound.Code, data.Code)
		assert.Equal(t, ErrorResponseLeaderboardNotFound.Message, data.Message)
	})

	t.Run("Invalid Leaderboard ID", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: leaderboard.BuildGetByIDAndGameIDFunc(func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{}, leaderboard.ErrInvalidLeaderboardID
			}),
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/leaderboards/%s", expectedID), nil)

		req.Header.Set(gameIDHeader, expectedGameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseLeaderboardInvalidID.Code, data.Code)
		assert.Equal(t, ErrorResponseLeaderboardInvalidID.Message, data.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		app := App(Config{
			GetLeaderboardByIDAndGameIDFunc: leaderboard.BuildGetByIDAndGameIDFunc(func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{}, errors.New("any error")
			}),
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/leaderboards/%s", expectedID), nil)

		req.Header.Set(gameIDHeader, expectedGameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseInternalServerError.Code, data.Code)
		assert.Equal(t, ErrorResponseInternalServerError.Message, data.Message)
	})
}

func TestBuildDeleteLeaderboardHandler(t *testing.T) {
	var (
		expectedID     = uuid.NewString()
		expectedGameID = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			DeleteLeaderboardByIDAndGameIDFunc: leaderboard.BuildSoftDeleteFunc(func(ctx context.Context, id, gameID string) error {
				return nil
			}),
		})

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/leaderboards/%s", expectedID), nil)

		req.Header.Set(gameIDHeader, expectedGameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Missing Game ID", func(t *testing.T) {
		app := App(Config{
			DeleteLeaderboardByIDAndGameIDFunc: leaderboard.BuildSoftDeleteFunc(func(ctx context.Context, id, gameID string) error {
				return nil
			}),
		})

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/leaderboards/%s", expectedID), nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseLeaderboardInvalidGameID.Code, data.Code)
		assert.Equal(t, ErrorResponseLeaderboardInvalidGameID.Message, data.Message)
	})

	t.Run("Not Found", func(t *testing.T) {
		app := App(Config{
			DeleteLeaderboardByIDAndGameIDFunc: leaderboard.BuildSoftDeleteFunc(func(ctx context.Context, id, gameID string) error {
				return leaderboard.ErrLeaderboardNotFound
			}),
		})

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/leaderboards/%s", expectedID), nil)

		req.Header.Set(gameIDHeader, expectedGameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseLeaderboardNotFound.Code, data.Code)
		assert.Equal(t, ErrorResponseLeaderboardNotFound.Message, data.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		app := App(Config{
			DeleteLeaderboardByIDAndGameIDFunc: leaderboard.BuildSoftDeleteFunc(func(ctx context.Context, id, gameID string) error {
				return errors.New("any error")
			}),
		})

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/leaderboards/%s", expectedID), nil)

		req.Header.Set(gameIDHeader, expectedGameID)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseInternalServerError.Code, data.Code)
		assert.Equal(t, ErrorResponseInternalServerError.Message, data.Message)
	})
}
