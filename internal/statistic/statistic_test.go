package statistic

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStatisticValidete(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		err := NewStatisticData{
			GameID:          uuid.NewString(),
			Name:            "Test Validate Statistic",
			Description:     "Test validate statistic unit test",
			AggregationMode: AggregationModeInc,
			Goal:            nil,
			Landmarks:       []float64{50, 100, 200},
		}.validate()

		assert.NoError(t, err)
	})

	t.Run("Validation Error", func(t *testing.T) {
		err := NewStatisticData{}.validate()

		assert.ErrorIs(t, err, ErrStatisticValidation)
		assert.ErrorIs(t, err, ErrMissingGameID)
		assert.ErrorIs(t, err, ErrInvalidName)
		assert.ErrorIs(t, err, ErrInvalidAggregationMode)
	})
}

func TestBuildCreateStatisticFunc(t *testing.T) {
	var (
		ctx = context.Background()

		statisticID = uuid.NewString()
		gameID      = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		createStatisticFunc := BuildCreateStatisticFunc(func(ctx context.Context, data NewStatisticData) (Statistic, error) {
			return Statistic{
				ID:              statisticID,
				GameID:          data.GameID,
				Name:            data.Name,
				Description:     data.Description,
				AggregationMode: data.AggregationMode,
				Goal:            data.Goal,
				Landmarks:       data.Landmarks,
			}, nil
		})

		data := NewStatisticData{
			GameID:          gameID,
			Name:            "Test Create Statistic",
			Description:     "Test build create statistic unit test",
			AggregationMode: AggregationModeMax,
			Goal:            nil,
			Landmarks:       []float64{10, 50, 100},
		}

		statistic, err := createStatisticFunc(ctx, data)
		assert.NoError(t, err)

		assert.Equal(t, statisticID, statistic.ID)
		assert.Equal(t, data.GameID, statistic.GameID)
		assert.Equal(t, data.Name, statistic.Name)
		assert.Equal(t, data.Description, statistic.Description)
		assert.Equal(t, data.AggregationMode, statistic.AggregationMode)
		assert.Equal(t, data.Goal, statistic.Goal)
		assert.Equal(t, data.Landmarks, statistic.Landmarks)
	})

	t.Run("Validation Error", func(t *testing.T) {
		var (
			createStatisticFunc = BuildCreateStatisticFunc(nil)
			data                = NewStatisticData{}
		)

		statistic, err := createStatisticFunc(ctx, data)

		assert.ErrorIs(t, err, ErrStatisticValidation)
		assert.Empty(t, statistic.ID)
	})

	t.Run("Random Error", func(t *testing.T) {
		createStatisticFunc := BuildCreateStatisticFunc(func(ctx context.Context, data NewStatisticData) (Statistic, error) {
			return Statistic{}, errors.New("any error")
		})

		data := NewStatisticData{
			GameID:          gameID,
			Name:            "Test Create Statistic",
			Description:     "Test build create statistic unit test",
			AggregationMode: AggregationModeMax,
			Goal:            nil,
			Landmarks:       []float64{10, 50, 100},
		}

		statistic, err := createStatisticFunc(ctx, data)

		assert.Error(t, err)
		assert.Empty(t, statistic.ID)
	})
}
