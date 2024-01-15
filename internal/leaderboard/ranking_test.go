package leaderboard

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildUpsertPlayerRankFunc(t *testing.T) {
	var (
		ctx = context.Background()

		leaderboardID = uuid.NewString()
		gameID        = uuid.NewString()
		playerID      = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		lb := Leaderboard{
			ID:              leaderboardID,
			GameID:          gameID,
			AggregationMode: AggregationModeInc,
		}

		upsertPlayerRankFunc := BuildUpsertPlayerRankFunc(func(ctx context.Context, leaderboard Leaderboard, playerID string, value float64) error {
			return nil
		})

		err := upsertPlayerRankFunc(ctx, lb, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("Invalid Aggregation Mode", func(t *testing.T) {
		lb := Leaderboard{
			ID:              leaderboardID,
			GameID:          gameID,
			AggregationMode: "INVALID",
		}

		upsertPlayerRankFunc := BuildUpsertPlayerRankFunc(func(ctx context.Context, leaderboard Leaderboard, playerID string, value float64) error {
			return ErrInvalidAggregationMode
		})

		err := upsertPlayerRankFunc(ctx, lb, playerID, rand.Float64())
		assert.ErrorIs(t, err, ErrInvalidAggregationMode)
	})

	t.Run("Leaderboard Closed", func(t *testing.T) {
		lb := Leaderboard{EndAt: time.Now().Add(-24 * time.Hour)}

		upsertPlayerRankFunc := BuildUpsertPlayerRankFunc(nil)

		err := upsertPlayerRankFunc(ctx, lb, playerID, rand.Float64())
		assert.ErrorIs(t, err, ErrLeaderboardClosed)
	})
}

func TestBuildRankingFunc(t *testing.T) {
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		lb := Leaderboard{
			ID:       uuid.NewString(),
			Ordering: OrderingAsc,
		}

		rankingFunc := BuildRankingFunc(func(ctx context.Context, leaderboardID, ordering string, page, limit int64) ([]Rank, error) {
			return make([]Rank, 0), nil
		})

		_, err := rankingFunc(ctx, lb, 0, 10)
		assert.NoError(t, err)
	})

	t.Run("Page Number Lower Than Minimun", func(t *testing.T) {
		lb := Leaderboard{
			ID:       uuid.NewString(),
			Ordering: OrderingAsc,
		}

		rankingFunc := BuildRankingFunc(nil)

		_, err := rankingFunc(ctx, lb, MinPageNumber-1, 10)
		assert.ErrorIs(t, err, ErrInvalidPageNumber)
	})

	t.Run("Limit Number Lower Than Minimun", func(t *testing.T) {
		lb := Leaderboard{
			ID:       uuid.NewString(),
			Ordering: OrderingAsc,
		}

		rankingFunc := BuildRankingFunc(nil)

		_, err := rankingFunc(ctx, lb, 0, MinLimitNumber-1)
		assert.ErrorIs(t, err, ErrInvalidLimitNumber)
	})

	t.Run("Limit Number Greater Than Maximum", func(t *testing.T) {
		lb := Leaderboard{
			ID:       uuid.NewString(),
			Ordering: OrderingAsc,
		}

		rankingFunc := BuildRankingFunc(nil)

		_, err := rankingFunc(ctx, lb, 0, MaxLimitNumber+1)
		assert.ErrorIs(t, err, ErrInvalidLimitNumber)
	})

	t.Run("Invalid Ordering Value", func(t *testing.T) {
		lb := Leaderboard{
			ID:       uuid.NewString(),
			Ordering: "INVALID",
		}

		rankingFunc := BuildRankingFunc(func(ctx context.Context, leaderboardID, ordering string, page, limit int64) ([]Rank, error) {
			return nil, ErrInvalidOrdering
		})

		_, err := rankingFunc(ctx, lb, 0, 10)
		assert.ErrorIs(t, err, ErrInvalidOrdering)
	})

	t.Run("Random Error", func(t *testing.T) {
		lb := Leaderboard{
			ID:       uuid.NewString(),
			Ordering: OrderingAsc,
		}

		rankingFunc := BuildRankingFunc(func(ctx context.Context, leaderboardID, ordering string, page, limit int64) ([]Rank, error) {
			return nil, errors.New("any error")
		})

		_, err := rankingFunc(ctx, lb, 0, 10)
		assert.Error(t, err)
	})
}
