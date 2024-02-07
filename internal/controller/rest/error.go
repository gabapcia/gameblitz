package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gabapcia/gameblitz/internal/auth"
	"github.com/gabapcia/gameblitz/internal/infra/logger/zap"
	"github.com/gabapcia/gameblitz/internal/leaderboard"
	"github.com/gabapcia/gameblitz/internal/quest"
	"github.com/gabapcia/gameblitz/internal/statistic"

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
		// Auth
		case errors.Is(err, auth.ErrInvalidCredentials):
			validationErrorMessages := strings.Split(err.Error(), "\n")
			return c.Status(http.StatusForbidden).JSON(ErrorResponseInvalidAuthCredentials.withDetails(validationErrorMessages...))
		// Statistic
		case errors.Is(err, statistic.ErrPlayerStatisticNotFound):
			return c.Status(http.StatusNotFound).JSON(ErrorResponsePlayerStatisticNotFound)
		case errors.Is(err, statistic.ErrInvalidStatisticID):
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseStatisticInvalidID)
		case errors.Is(err, statistic.ErrStatisticNotFound):
			return c.Status(http.StatusNotFound).JSON(ErrorResponseStatisticNotFound)
		case errors.Is(err, statistic.ErrStatisticValidation):
			validationErrorMessages := strings.Split(err.Error(), "\n")
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseStatisticInvalid.withDetails(validationErrorMessages...))
		// Quest
		case errors.Is(err, quest.ErrInvalidQuestID):
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseQuestInvalidID)
		case errors.Is(err, quest.ErrQuestNotFound):
			return c.Status(http.StatusNotFound).JSON(ErrorResponseQuestNotFound)
		case errors.Is(err, quest.ErrPlayerAlreadyStartedTheQuest):
			return c.Status(http.StatusConflict).JSON(ErrorResponsePlayerAlreadyStartedTheQuest)
		case errors.Is(err, quest.ErrPlayerNotStartedTheQuest):
			return c.Status(http.StatusNotFound).JSON(ErrorResponsePlayerNotStartedTheQuest)
		case errors.Is(err, quest.ErrPlayerQuestAlreadyCompleted):
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponsePlayerQuestAlreadyFinished)
		case errors.Is(err, quest.ErrQuestValidationError):
			validationErrorMessages := strings.Split(err.Error(), "\n")
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseQuestInvalid.withDetails(validationErrorMessages...))
			// Leaderboard
		case errors.Is(err, leaderboard.ErrLeaderboardClosed):
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseLeaderboardClosed)
		case errors.Is(err, leaderboard.ErrInvalidPageNumber):
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseRankingPageNumber)
		case errors.Is(err, leaderboard.ErrInvalidLimitNumber):
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseRankingLimitNumber)
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
