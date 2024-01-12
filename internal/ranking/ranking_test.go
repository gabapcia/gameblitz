package ranking

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"
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

	t.Run("Increment By Value", func(t *testing.T) {
		upsertPlayerRankFunc := BuildUpsertPlayerRankFunc(
			func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{
					ID:              id,
					GameID:          gameID,
					AggregationMode: leaderboard.AggregationModeInc,
				}, nil
			},
			func(ctx context.Context, leaderboardID, playerID string, value float64) error {
				return nil
			},
			nil,
			nil,
		)

		err := upsertPlayerRankFunc(ctx, leaderboardID, gameID, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("Set Max Value", func(t *testing.T) {
		upsertPlayerRankFunc := BuildUpsertPlayerRankFunc(
			func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{
					ID:              id,
					GameID:          gameID,
					AggregationMode: leaderboard.AggregationModeMax,
				}, nil
			},
			nil,
			func(ctx context.Context, leaderboardID, playerID string, value float64) error {
				return nil
			},
			nil,
		)

		err := upsertPlayerRankFunc(ctx, leaderboardID, gameID, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("Set Min Value", func(t *testing.T) {
		upsertPlayerRankFunc := BuildUpsertPlayerRankFunc(
			func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{
					ID:              id,
					GameID:          gameID,
					AggregationMode: leaderboard.AggregationModeMin,
				}, nil
			},
			nil,
			nil,
			func(ctx context.Context, leaderboardID, playerID string, value float64) error {
				return nil
			},
		)

		err := upsertPlayerRankFunc(ctx, leaderboardID, gameID, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("Invalid Aggregation Mode", func(t *testing.T) {
		upsertPlayerRankFunc := BuildUpsertPlayerRankFunc(
			func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{
					ID:              id,
					GameID:          gameID,
					AggregationMode: "INVALID",
				}, nil
			},
			nil,
			nil,
			nil,
		)

		err := upsertPlayerRankFunc(ctx, leaderboardID, gameID, playerID, rand.Float64())
		assert.ErrorIs(t, err, ErrInvalidAggregationMode)
	})

	t.Run("Error On Get Leaderboard", func(t *testing.T) {
		upsertPlayerRankFunc := BuildUpsertPlayerRankFunc(
			func(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
				return leaderboard.Leaderboard{}, errors.New("any error")
			},
			nil,
			nil,
			nil,
		)

		err := upsertPlayerRankFunc(ctx, leaderboardID, gameID, playerID, rand.Float64())
		assert.Error(t, err)
	})
}
