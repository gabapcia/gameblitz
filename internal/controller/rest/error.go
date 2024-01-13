package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gabarcia/metagaming-api/internal/infra/logger/zap"
	"github.com/gabarcia/metagaming-api/internal/leaderboard"
	"github.com/gabarcia/metagaming-api/internal/ranking"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Code    string   `json:"code"`              // Error unique code
	Message string   `json:"message"`           // Error message
	Details []string `json:"details,omitempty"` // Details about the source of the error
}

func (e ErrorResponse) withDetails(details ...string) ErrorResponse {
	e.Details = details
	return e
}

var (
	ErrorResponseInternalServerError = ErrorResponse{Code: "0.0", Message: "Unknown error"}
	ErrorResponseInvalidRequestBody  = ErrorResponse{Code: "0.1", Message: "Invalid request body"}
)

func buildErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var jsonErr *json.SyntaxError

		switch {
		// Ranking
		case errors.Is(err, ranking.ErrLeaderboardClosed):
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseLeaderboardClosed)
		case errors.Is(err, ranking.ErrInvalidPageNumber):
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseRankingPageNumber)
		case errors.Is(err, ranking.ErrInvalidLimitNumber):
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseRankingLimitNumber)
		// Leaderboard
		case errors.Is(err, leaderboard.ErrInvalidLeaderboardID):
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseLeaderboardInvalidID)
		case errors.Is(err, leaderboard.ErrLeaderboardNotFound):
			return c.Status(http.StatusNotFound).JSON(ErrorResponseLeaderboardNotFound)
		case errors.Is(err, leaderboard.ErrValidationError):
			validationErrorMessages := strings.Split(err.Error(), "\n")
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseLeaderboardInvalid.withDetails(validationErrorMessages...))
		// Unknown
		case errors.As(err, &jsonErr):
			return c.Status(http.StatusBadRequest).JSON(ErrorResponseInvalidRequestBody)
		default:
			zap.Error(err, "unknown error")
			return c.Status(http.StatusInternalServerError).JSON(ErrorResponseInternalServerError)
		}
	}
}
