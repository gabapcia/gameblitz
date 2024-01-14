package statistic

import (
	"context"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildUpdatePlayerProgressionFunc(t *testing.T) {
	var (
		ctx = context.Background()

		playerID = uuid.NewString()
	)

	t.Run("OK Without Goals Or Landmarks And With Aggregation Mode INC", func(t *testing.T) {
		updatePlayerProgressionFunc := BuildUpdatePlayerProgressionFunc(
			nil,
			func(ctx context.Context, statisticID, gameID, playerID string, value float64) (PlayerProgressionUpdates, error) {
				return PlayerProgressionUpdates{
					PlayerID:           playerID,
					CurrentValue:       value,
					GoalComplted:       false,
					LandmarksCompleted: make([]float64, 0),
				}, nil
			},
			nil,
			nil,
		)

		statistic := Statistic{
			AggregationMode: AggregationModeInc,
		}

		err := updatePlayerProgressionFunc(ctx, statistic, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("OK Without Goals Or Landmarks And With Aggregation Mode MAX", func(t *testing.T) {
		updatePlayerProgressionFunc := BuildUpdatePlayerProgressionFunc(
			nil,
			nil,
			func(ctx context.Context, statisticID, gameID, playerID string, value float64) (PlayerProgressionUpdates, error) {
				return PlayerProgressionUpdates{
					PlayerID:           playerID,
					CurrentValue:       value,
					GoalComplted:       false,
					LandmarksCompleted: make([]float64, 0),
				}, nil
			},
			nil,
		)

		statistic := Statistic{
			AggregationMode: AggregationModeMax,
		}

		err := updatePlayerProgressionFunc(ctx, statistic, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("OK Without Goals Or Landmarks And With Aggregation Mode MIN", func(t *testing.T) {
		updatePlayerProgressionFunc := BuildUpdatePlayerProgressionFunc(
			nil,
			nil,
			nil,
			func(ctx context.Context, statisticID, gameID, playerID string, value float64) (PlayerProgressionUpdates, error) {
				return PlayerProgressionUpdates{
					PlayerID:           playerID,
					CurrentValue:       value,
					GoalComplted:       false,
					LandmarksCompleted: make([]float64, 0),
				}, nil
			},
		)

		statistic := Statistic{
			AggregationMode: AggregationModeMin,
		}

		err := updatePlayerProgressionFunc(ctx, statistic, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("OK With Goals Or Landmarks", func(t *testing.T) {
		updatePlayerProgressionFunc := BuildUpdatePlayerProgressionFunc(
			func(ctx context.Context, statistic Statistic, updated PlayerProgressionUpdates) error {
				return nil
			},
			func(ctx context.Context, statisticID, gameID, playerID string, value float64) (PlayerProgressionUpdates, error) {
				return PlayerProgressionUpdates{
					PlayerID:           playerID,
					CurrentValue:       value,
					GoalComplted:       true,
					LandmarksCompleted: []float64{10, 50},
				}, nil
			},
			nil,
			nil,
		)

		statistic := Statistic{
			AggregationMode: AggregationModeInc,
		}

		err := updatePlayerProgressionFunc(ctx, statistic, playerID, rand.Float64())
		assert.NoError(t, err)
	})

	t.Run("Invalid Aggregation Mode", func(t *testing.T) {
		var (
			updatePlayerProgressionFunc = BuildUpdatePlayerProgressionFunc(nil, nil, nil, nil)
			statistic                   = Statistic{AggregationMode: "INVALID"}
		)

		err := updatePlayerProgressionFunc(ctx, statistic, playerID, rand.Float64())
		assert.ErrorIs(t, err, ErrInvalidAggregationMode)
	})
}
