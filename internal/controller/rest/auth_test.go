package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabapcia/gameblitz/internal/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildAuthMiddleware(t *testing.T) {
	var (
		gameID = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		authMiddleware := buildAuthMiddleware(func(ctx context.Context, credentials string) (auth.Claims, error) {
			return auth.Claims{GameID: gameID}, nil
		})

		app := fiber.New()
		app.Get("/", authMiddleware, func(c *fiber.Ctx) error {
			claims := c.Locals("claims").(auth.Claims)
			return c.Status(http.StatusOK).JSON(claims)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body auth.Claims
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, gameID, body.GameID)
	})

	t.Run("Missing Token", func(t *testing.T) {
		authMiddleware := buildAuthMiddleware(nil)

		app := fiber.New()
		app.Get("/", authMiddleware, func(c *fiber.Ctx) error {
			claims := c.Locals("claims").(auth.Claims)
			return c.Status(http.StatusOK).JSON(claims)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseMissingAuthCredentials.Code, body.Code)
		assert.Equal(t, ErrorResponseMissingAuthCredentials.Message, body.Message)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		authMiddleware := buildAuthMiddleware(func(ctx context.Context, credentials string) (auth.Claims, error) {
			return auth.Claims{}, auth.ErrInvalidCredentials
		})

		app := fiber.New(fiber.Config{ErrorHandler: buildErrorHandler()})
		app.Get("/", authMiddleware, func(c *fiber.Ctx) error {
			claims := c.Locals("claims").(auth.Claims)
			return c.Status(http.StatusOK).JSON(claims)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		req.Header.Set("Authorization", uuid.NewString())

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		var body ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, ErrorResponseInvalidAuthCredentials.Code, body.Code)
		assert.Equal(t, ErrorResponseInvalidAuthCredentials.Message, body.Message)
	})
}
