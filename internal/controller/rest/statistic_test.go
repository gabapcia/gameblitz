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

	"github.com/gabarcia/gameblitz/internal/auth"
	"github.com/gabarcia/gameblitz/internal/infra/logger/zap"
	"github.com/gabarcia/gameblitz/internal/statistic"
	"github.com/gofiber/fiber/v2"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildGetStatisticMiddleware(t *testing.T) {
	var (
		leaderboardID = uuid.NewString()
		gameID        = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		authMiddleware := buildAuthMiddleware(func(ctx context.Context, credentials string) (auth.Claims, error) {
			return auth.Claims{GameID: gameID}, nil
		})
		getStatisticMiddleware := buildGetStatisticMiddleware(nil, time.Minute, func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
			return statistic.Statistic{ID: id, GameID: gameID}, nil
		})

		app := fiber.New()
		app.Get("/:statisticId", authMiddleware, getStatisticMiddleware, func(c *fiber.Ctx) error {
			statistic := c.Locals("statistic")
			return c.Status(http.StatusOK).JSON(statistic)
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", leaderboardID), nil)

		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body Leaderboard
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, leaderboardID, body.ID)
		assert.Equal(t, gameID, body.GameID)
	})

	t.Run("Leaderboard Not Found", func(t *testing.T) {
		authMiddleware := buildAuthMiddleware(func(ctx context.Context, credentials string) (auth.Claims, error) {
			return auth.Claims{GameID: gameID}, nil
		})
		getStatisticMiddleware := buildGetStatisticMiddleware(nil, time.Minute, func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
			return statistic.Statistic{}, statistic.ErrStatisticNotFound
		})

		app := fiber.New(fiber.Config{ErrorHandler: buildErrorHandler()})
		app.Get("/:statisticId", authMiddleware, getStatisticMiddleware, func(c *fiber.Ctx) error {
			statistic := c.Locals("statistic")
			return c.Status(http.StatusOK).JSON(statistic)
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", leaderboardID), nil)

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
}

func TestBuildCreateStatisticHanlder(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		gameID := uuid.NewString()
		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			CreateStatisticFunc: func(ctx context.Context, data statistic.NewStatisticData) (statistic.Statistic, error) {
				return statistic.Statistic{
					ID:              uuid.NewString(),
					GameID:          data.GameID,
					Name:            data.Name,
					Description:     data.Description,
					AggregationMode: data.AggregationMode,
					Goal:            data.Goal,
					Landmarks:       data.Landmarks,
				}, nil
			},
		})

		body, err := json.Marshal(map[string]any{
			"name":            "Test Create Statistic",
			"description":     "Test create statistic handler unit test",
			"aggregationMode": "MAX",
			"canOverflow":     true,
			"goal":            nil,
			"landmarks":       []float64{10, 50, 100},
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/statistics", bytes.NewBuffer(body))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var data Statistic
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.NotEmpty(t, data.ID)
		assert.Equal(t, gameID, data.GameID)
	})

	t.Run("Validation Error", func(t *testing.T) {
		gameID := uuid.NewString()
		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			CreateStatisticFunc: statistic.BuildCreateStatisticFunc(nil),
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/statistics", bytes.NewBufferString("{}"))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseStatisticInvalid.Code, data.Code)
		assert.Equal(t, ErrorResponseStatisticInvalid.Message, data.Message)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		var (
			gameID = uuid.NewString()
			app    = App(Config{
				AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
					return auth.Claims{GameID: gameID}, nil
				},
			})
			req = httptest.NewRequest(http.MethodPost, "/api/v1/statistics", bytes.NewBufferString("{"))
		)

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

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		gameID := uuid.NewString()
		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			CreateStatisticFunc: func(ctx context.Context, data statistic.NewStatisticData) (statistic.Statistic, error) {
				return statistic.Statistic{}, errors.New("any error")
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/statistics", bytes.NewBufferString("{}"))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

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

func TestBuildGetStatisticHanlder(t *testing.T) {
	var (
		statisticID = uuid.NewString()
		gameID      = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{ID: id, GameID: gameID}, nil
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/statistics/%s", statisticID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var data Statistic
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, statisticID, data.ID)
		assert.Equal(t, gameID, data.GameID)
	})

	t.Run("Invalid Statistic ID", func(t *testing.T) {
		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{}, statistic.ErrInvalidStatisticID
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/statistics/%s", statisticID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseStatisticInvalidID.Code, data.Code)
		assert.Equal(t, ErrorResponseStatisticInvalidID.Message, data.Message)
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

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/statistics/%s", statisticID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseStatisticNotFound.Code, data.Code)
		assert.Equal(t, ErrorResponseStatisticNotFound.Message, data.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			GetStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
				return statistic.Statistic{}, errors.New("any error")
			},
		})

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/statistics/%s", statisticID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

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

func TestBuildDeleteStatisticHanlder(t *testing.T) {
	var (
		statisticID = uuid.NewString()
		gameID      = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			SoftDeleteStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) error {
				return nil
			},
		})

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/statistics/%s", statisticID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Invalid Statistic ID", func(t *testing.T) {
		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			SoftDeleteStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) error {
				return statistic.ErrInvalidStatisticID
			},
		})

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/statistics/invalid-id", nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseStatisticInvalidID.Code, data.Code)
		assert.Equal(t, ErrorResponseStatisticInvalidID.Message, data.Message)
	})

	t.Run("Not Found", func(t *testing.T) {
		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			SoftDeleteStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) error {
				return statistic.ErrStatisticNotFound
			},
		})

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/statistics/%s", statisticID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var data ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseStatisticNotFound.Code, data.Code)
		assert.Equal(t, ErrorResponseStatisticNotFound.Message, data.Message)
	})

	t.Run("Random Error", func(t *testing.T) {
		zap.Start()
		defer zap.Sync()

		app := App(Config{
			AuthenticateFunc: func(ctx context.Context, credentials string) (auth.Claims, error) {
				return auth.Claims{GameID: gameID}, nil
			},
			SoftDeleteStatisticByIDAndGameIDFunc: func(ctx context.Context, id, gameID string) error {
				return errors.New("any error")
			},
		})

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/statistics/%s", statisticID), nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", uuid.NewString())

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
