package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gabarcia/metagaming-api/internal/infra/logger/zap"
	"github.com/gabarcia/metagaming-api/internal/leaderboard"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func (e ErrorResponse) withDetails(details ...string) ErrorResponse {
	e.Details = details
	return e
}

var (
	ErrorResponseInternalServerError = ErrorResponse{Code: "0.1", Message: "Unknown error"}
)

func BuildErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		switch {
		case errors.Is(err, leaderboard.ErrInvalidLeaderboardID):
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseLeaderboardInvalidID)
		case errors.Is(err, leaderboard.ErrLeaderboardNotFound):
			return c.Status(http.StatusNotFound).JSON(ErrorResponseLeaderboardNotFound)
		case errors.Is(err, leaderboard.ErrValidationError):
			validationErrorMessages := strings.Split(err.Error(), "\n")
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseLeaderboardInvalid.withDetails(validationErrorMessages...))
		default:
			zap.Error(err, "unknown error")
			return c.Status(http.StatusInternalServerError).JSON(ErrorResponseInternalServerError)
		}
	}
}
