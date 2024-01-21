package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gabarcia/gameblitz/internal/auth"
	"github.com/gabarcia/gameblitz/internal/infra/logger/zap"
	"github.com/gabarcia/gameblitz/internal/leaderboard"

	"github.com/gofiber/fiber/v2"
)

type CreateLeaderboardReq struct {
	Name            string    `json:"name"`                                // Leaderboard's name
	Description     string    `json:"description"`                         // Leaderboard's description
	StartAt         time.Time `json:"startAt"`                             // Time that the leaderboard should start working
	EndAt           time.Time `json:"endAt"`                               // Time that the leaderboard will be closed for new updates
	AggregationMode string    `json:"aggregationMode" enums:"INC,MAX,MIN"` // Data aggregation mode
	Ordering        string    `json:"ordering" enums:"ASC,DESC"`           // Leaderboard ranking order
}

type Leaderboard struct {
	CreatedAt       time.Time  `json:"createdAt"`                           // Time that the leaderboard was created
	UpdatedAt       time.Time  `json:"updatedAt"`                           // Last time that the leaderboard info was updated
	ID              string     `json:"id"`                                  // Leaderboard's ID
	GameID          string     `json:"gameId"`                              // The ID from the game that is responsible for the leaderboard
	Name            string     `json:"name"`                                // Leaderboard's name
	Description     string     `json:"description"`                         // Leaderboard's description
	StartAt         time.Time  `json:"startAt"`                             // Time that the leaderboard should start working
	EndAt           *time.Time `json:"endAt"`                               // Time that the leaderboard will be closed for new updates
	AggregationMode string     `json:"aggregationMode" enums:"INC,MAX,MIN"` // Data aggregation mode
	Ordering        string     `json:"ordering" enums:"ASC,DESC"`           // Leaderboard ranking order
}

func (r CreateLeaderboardReq) toDomain(gameID string) leaderboard.NewLeaderboardData {
	return leaderboard.NewLeaderboardData{
		GameID:          gameID,
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
	ErrorResponseLeaderboardInvalid   = ErrorResponse{Code: "1.0", Message: "Invalid leaderboard"}
	ErrorResponseLeaderboardNotFound  = ErrorResponse{Code: "1.1", Message: "Leaderboard not found"}
	ErrorResponseLeaderboardInvalidID = ErrorResponse{Code: "1.2", Message: "Invalid leaderboard ID"}
)

func buildGetLeaderboardMiddleware(cache fiber.Storage, expiration time.Duration, getLeaderboardByIDAndGameIDFunc leaderboard.GetByIDAndGameIDFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			id       = c.Params("leaderboardId")
			claims   = c.Locals("claims").(auth.Claims)
			cacheKey = fmt.Sprintf("GetLeaderboardMiddleware:%s:%s", id, claims.GameID)
		)

		if cache != nil {
			data, err := cache.Get(cacheKey)
			if err != nil {
				zap.Error(err, "get cache error")
			} else if data != nil {
				var leaderboard leaderboard.Leaderboard
				if err = json.Unmarshal(data, &leaderboard); err != nil {
					zap.Error(err, "unmarshal cached leaderboard error")
				} else {
					c.Locals("leaderboard", leaderboard)
					return c.Next()
				}
			}
		}

		leaderboard, err := getLeaderboardByIDAndGameIDFunc(c.Context(), id, claims.GameID)
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

// @summary Create Leaderboard
// @description Create a leaderboard
// @router /api/v1/leaderboards [POST]
// @accept json
// @produce json
// @param Authorization header string true "Game's JWT authorization"
// @param NewLeaderboardData body CreateLeaderboardReq true "New leaderboard config data"
// @success 201 {object} Leaderboard
// @failure 400,422,500 {object} ErrorResponse
func buildCreateLeaderboardHandler(createLeaderboardFunc leaderboard.CreateFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(auth.Claims)

		var body CreateLeaderboardReq
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		leaderboard, err := createLeaderboardFunc(c.Context(), body.toDomain(claims.GameID))
		if err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(leaderboardFromDomain(leaderboard))
	}
}

// @summary Get Leaderboard
// @description Return a leaderboard by id and game id
// @router /api/v1/leaderboards/{leaderboardId} [GET]
// @produce json
// @param Authorization header string true "Game's JWT authorization"
// @param leaderboardId path string true "Leaderboard ID"
// @success 200 {object} Leaderboard
// @failure 404,422,500 {object} ErrorResponse
func buildGetLeaderboardHandler(getLeaderboardByIDAndGameIDFunc leaderboard.GetByIDAndGameIDFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			id     = c.Params("leaderboardId")
			claims = c.Locals("claims").(auth.Claims)
		)

		leaderboard, err := getLeaderboardByIDAndGameIDFunc(c.Context(), id, claims.GameID)
		if err != nil {
			return err
		}

		return c.Status(http.StatusOK).JSON(leaderboardFromDomain(leaderboard))
	}
}

// @summary Delete Leaderboard
// @description Delete a leaderboard by id and game id
// @router /api/v1/leaderboards/{leaderboardId} [DELETE]
// @param Authorization header string true "Game's JWT authorization"
// @param leaderboardId path string true "Leaderboard ID"
// @success 204
// @failure 404,422,500 {object} ErrorResponse
func buildDeleteLeaderboardHandler(deleteLeaderboardByIDAndGameIDFunc leaderboard.SoftDeleteFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			id     = c.Params("leaderboardId")
			claims = c.Locals("claims").(auth.Claims)
		)

		if err := deleteLeaderboardByIDAndGameIDFunc(c.Context(), id, claims.GameID); err != nil {
			return err
		}

		return c.SendStatus(http.StatusNoContent)
	}
}
