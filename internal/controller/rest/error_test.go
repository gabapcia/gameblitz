package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabarcia/game-blitz/internal/infra/logger/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponseWithDetails(t *testing.T) {
	var (
		testCode    = "Test Code"
		testMessage = "Test Message"

		expected = ErrorResponse{Code: testCode, Message: testMessage}

		details = []string{"detail 1", "detail 2"}
	)

	assert.Empty(t, expected.Details)

	got := expected.withDetails(details...)

	assert.Equal(t, testCode, got.Code)
	assert.Equal(t, testMessage, got.Message)
	assert.Equal(t, details, got.Details)
}

func TestBuildErrorHandler(t *testing.T) {
	zap.Start()
	defer zap.Sync()

	app := fiber.New(fiber.Config{ErrorHandler: buildErrorHandler()})
	app.Get("/", func(c *fiber.Ctx) error {
		return errors.New("any error")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var body ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)

	assert.Equal(t, ErrorResponseInternalServerError.Code, body.Code)
	assert.Equal(t, ErrorResponseInternalServerError.Message, body.Message)
}
