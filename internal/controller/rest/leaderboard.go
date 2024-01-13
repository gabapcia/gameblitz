package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gabarcia/metagaming-api/internal/infra/logger/zap"
	"github.com/gabarcia/metagaming-api/internal/leaderboard"

	"github.com/gofiber/fiber/v2"
)

type CreateLeaderboardReq struct {
	GameID          string    `json:"gameId"`          // The ID from the game that is responsible for the leaderboard
	Name            string    `json:"name"`            // Leaderboard's name
	Description     string    `json:"description"`     // Leaderboard's description
	StartAt         time.Time `json:"startAt"`         // Time that the leaderboard should start working
	EndAt           time.Time `json:"endAt"`           // Time that the leaderboard will be closed for new updates
	AggregationMode string    `json:"aggregationMode"` // Data aggregation mode
	Ordering        string    `json:"ordering"`        // Leaderboard ranking order
}

type Leaderboard struct {
	CreatedAt       time.Time  `json:"createdAt"`       // Time that the leaderboard was created
	UpdatedAt       time.Time  `json:"updatedAt"`       // Last time that the leaderboard info was updated
	ID              string     `json:"id"`              // Leaderboard's ID
	GameID          string     `json:"gameId"`          // The ID from the game that is responsible for the leaderboard
	Name            string     `json:"name"`            // Leaderboard's name
	Description     string     `json:"description"`     // Leaderboard's description
	StartAt         time.Time  `json:"startAt"`         // Time that the leaderboard should start working
	EndAt           *time.Time `json:"endAt"`           // Time that the leaderboard will be closed for new updates
	AggregationMode string     `json:"aggregationMode"` // Data aggregation mode
	Ordering        string     `json:"ordering"`        // Leaderboard ranking order
}

func (r CreateLeaderboardReq) toDomain() leaderboard.NewLeaderboardData {
	return leaderboard.NewLeaderboardData{
		GameID:          r.GameID,
		Name:            r.Name,
		Description:     r.Description,
		StartAt:         r.StartAt,
		EndAt:           r.EndAt,
		AggregationMode: r.AggregationMode,
		Ordering:        r.Ordering,
	}
}

func leaderboardFromDomain(l leaderboard.Leaderboard) Leaderboard {
	var endAt *time.Time
	if !l.EndAt.IsZero() {
		endAt = &l.EndAt
	}

	return Leaderboard{
		CreatedAt:       l.CreatedAt,
		UpdatedAt:       l.UpdatedAt,
		ID:              l.ID,
		GameID:          l.GameID,
		Name:            l.Name,
		Description:     l.Description,
		StartAt:         l.StartAt,
		EndAt:           endAt,
		AggregationMode: l.AggregationMode,
		Ordering:        l.Ordering,
	}
}

var (
	ErrorResponseLeaderboardInvalid       = ErrorResponse{Code: "1.0", Message: "Invalid Leaderboard"}
	ErrorResponseLeaderboardNotFound      = ErrorResponse{Code: "1.1", Message: "Leaderboard not found"}
	ErrorResponseLeaderboardInvalidID     = ErrorResponse{Code: "1.2", Message: "Invalid Leaderboard ID"}
	ErrorResponseLeaderboardInvalidGameID = ErrorResponse{Code: "1.3", Message: "Invalid Leaderboard Game ID"}
)

func buildGetLeaderboardMiddleware(cache fiber.Storage, expiration time.Duration, getLeaderboardByIDAndGameIDFunc leaderboard.GetByIDAndGameIDFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			id       = c.Params("leaderboardId")
			gameID   = string(c.Request().Header.Peek(gameIDHeader))
			cacheKey = fmt.Sprintf("GetLeaderboardMiddleware:%s:%s", id, gameID)
		)

		if gameID == "" {
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseLeaderboardInvalidGameID)
		}

		if cache != nil {
			data, err := cache.Get(cacheKey)
			if err != nil {
				zap.Error(err, "get cache error")
			} else {
				var leaderboard leaderboard.Leaderboard
				if err = json.Unmarshal(data, &leaderboard); err != nil {
					zap.Error(err, "unmarshal cached leaderboard error")
				} else {
					c.Locals("leaderboard", data)
					return c.Next()
				}
			}
		}

		leaderboard, err := getLeaderboardByIDAndGameIDFunc(c.Context(), id, gameID)
		if err != nil {
			return err
		}

		if cache != nil {
			data, err := json.Marshal(leaderboard)
			if err != nil {
				zap.Error(err, "marshal leaderboard cache error")
			} else {
				if err = cache.Set(cacheKey, data, expiration); err != nil {
					zap.Error(err, "unable to cache leaderboard")
				}
			}
		}

		c.Locals("leaderboard", leaderboard)
		return c.Next()
	}
}

func buildCreateLeaderboardHandler(createLeaderboardFunc leaderboard.CreateFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body CreateLeaderboardReq
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		leaderboard, err := createLeaderboardFunc(c.Context(), body.toDomain())
		if err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(leaderboardFromDomain(leaderboard))
	}
}

func buildGetLeaderboardHandler(getLeaderboardByIDAndGameIDFunc leaderboard.GetByIDAndGameIDFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			id     = c.Params("leaderboardId")
			gameID = string(c.Request().Header.Peek(gameIDHeader))
		)

		if gameID == "" {
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseLeaderboardInvalidGameID)
		}

		leaderboard, err := getLeaderboardByIDAndGameIDFunc(c.Context(), id, gameID)
		if err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(leaderboardFromDomain(leaderboard))
	}
}

func buildDeleteLeaderboardHandler(deleteLeaderboardByIDAndGameIDFunc leaderboard.SoftDeleteFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			id     = c.Params("leaderboardId")
			gameID = string(c.Request().Header.Peek(gameIDHeader))
		)

		if gameID == "" {
			return c.Status(http.StatusUnprocessableEntity).JSON(ErrorResponseLeaderboardInvalidGameID)
		}

		if err := deleteLeaderboardByIDAndGameIDFunc(c.Context(), id, gameID); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	}
}
