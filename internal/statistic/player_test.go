package statistic

import (
	"context"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildUpsertPlayerProgressionFunc(t *testing.T) {
	var (
		ctx = context.Background()

		playerID = uuid.NewString()
	)

	t.Run("OK Without Goals Or Landmarks And With Aggregation Mode SUM", func(t *testing.T) {
		updatePlayerProgressionFunc := BuildUpsertPlayerProgressionFunc(
			nil,
			func(ctx context.Context, statistic Statistic, playerID string, value float64) (PlayerProgression, PlayerProgressionUpdates, error) {
				return PlayerProgression{}, PlayerProgressionUpdates{}, nil
			},
		)

		statistic := Statistic{
			AggregationMode: AggregationModeSum,
		}

		err := updatePlayerProgressionFunc(ctx, statistic, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("OK Without Goals Or Landmarks And With Aggregation Mode SUB", func(t *testing.T) {
		updatePlayerProgressionFunc := BuildUpsertPlayerProgressionFunc(
			nil,
			func(ctx context.Context, statistic Statistic, playerID string, value float64) (PlayerProgression, PlayerProgressionUpdates, error) {
				return PlayerProgression{}, PlayerProgressionUpdates{}, nil
			},
		)

		statistic := Statistic{
			AggregationMode: AggregationModeSub,
		}

		err := updatePlayerProgressionFunc(ctx, statistic, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("OK Without Goals Or Landmarks And With Aggregation Mode MAX", func(t *testing.T) {
		updatePlayerProgressionFunc := BuildUpsertPlayerProgressionFunc(
			nil,
			func(ctx context.Context, statistic Statistic, playerID string, value float64) (PlayerProgression, PlayerProgressionUpdates, error) {
				return PlayerProgression{}, PlayerProgressionUpdates{}, nil
			},
		)

		statistic := Statistic{
			AggregationMode: AggregationModeMax,
		}

		err := updatePlayerProgressionFunc(ctx, statistic, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("OK Without Goals Or Landmarks And With Aggregation Mode MIN", func(t *testing.T) {
		updatePlayerProgressionFunc := BuildUpsertPlayerProgressionFunc(
			nil,
			func(ctx context.Context, statistic Statistic, playerID string, value float64) (PlayerProgression, PlayerProgressionUpdates, error) {
				return PlayerProgression{}, PlayerProgressionUpdates{}, nil
			},
		)

		statistic := Statistic{
			AggregationMode: AggregationModeMin,
		}

		err := updatePlayerProgressionFunc(ctx, statistic, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("OK With Goals Or Landmarks", func(t *testing.T) {
		updatePlayerProgressionFunc := BuildUpsertPlayerProgressionFunc(
			func(ctx context.Context, statistic Statistic, progression PlayerProgression, updates PlayerProgressionUpdates) error {
				return nil
			},
			func(ctx context.Context, statistic Statistic, playerID string, value float64) (PlayerProgression, PlayerProgressionUpdates, error) {
				return PlayerProgression{}, PlayerProgressionUpdates{}, nil
			},
		)

		statistic := Statistic{
			AggregationMode: AggregationModeSum,
		}

		err := updatePlayerProgressionFunc(ctx, statistic, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("Invalid Aggregation Mode", func(t *testing.T) {
		var (
			updatePlayerProgressionFunc = BuildUpsertPlayerProgressionFunc(
				nil,
				func(ctx context.Context, statistic Statistic, playerID string, value float64) (PlayerProgression, PlayerProgressionUpdates, error) {
					return PlayerProgression{}, PlayerProgressionUpdates{}, ErrInvalidAggregationMode
				},
			)

			statistic = Statistic{AggregationMode: "INVALID"}
		)

		err := updatePlayerProgressionFunc(ctx, statistic, playerID, rand.Float64())
		assert.ErrorIs(t, err, ErrInvalidAggregationMode)
	})
}
