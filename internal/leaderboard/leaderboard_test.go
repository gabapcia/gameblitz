package leaderboard

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		leaderboard := Leaderboard{
			GameID:          uuid.NewString(),
			Name:            "Test Leaderboard",
			Description:     "Test leaderboard validation unit test",
			StartAt:         time.Now(),
			EndAt:           time.Time{},
			AggregationMode: AggregationModeMax,
			DataType:        DataTypeInt,
			Ordering:        OrderingDesc,
		}

		assert.NoError(t, leaderboard.validate())
	})

	t.Run("Invalid Fields", func(t *testing.T) {
		leaderboard := Leaderboard{
			GameID:          "",
			Name:            "",
			Description:     "Test leaderboard validation unit test",
			StartAt:         time.Time{},
			EndAt:           time.Time{},
			AggregationMode: "INVALID",
			DataType:        "INVALID",
			Ordering:        "INVALID",
		}

		assert.ErrorIs(t, leaderboard.validate(), ErrValidationError)
		assert.ErrorIs(t, leaderboard.validate(), ErrInvalidGameID)
		assert.ErrorIs(t, leaderboard.validate(), ErrInvalidName)
		assert.ErrorIs(t, leaderboard.validate(), ErrInvalidStartDate)
		assert.ErrorIs(t, leaderboard.validate(), ErrInvalidAggregationMode)
		assert.ErrorIs(t, leaderboard.validate(), ErrInvalidDataType)
		assert.ErrorIs(t, leaderboard.validate(), ErrInvalidOrdering)
	})

	t.Run("End Date Before Start Date", func(t *testing.T) {
		leaderboard := Leaderboard{
			GameID:          uuid.NewString(),
			Name:            "Test Leaderboard",
			Description:     "Test leaderboard validation unit test",
			StartAt:         time.Now(),
			EndAt:           time.Now().Add(-24 * time.Hour),
			AggregationMode: AggregationModeMax,
			DataType:        DataTypeInt,
			Ordering:        OrderingDesc,
		}

		assert.ErrorIs(t, leaderboard.validate(), ErrValidationError)
		assert.ErrorIs(t, leaderboard.validate(), ErrEndDateBeforeStartDate)
	})
}

func TestBuildCreateFunc(t *testing.T) {
	var (
		ctx         = context.Background()
		expectedID  = uuid.NewString()
		leaderboard = Leaderboard{
			GameID:          uuid.NewString(),
			Name:            "Test Leaderboard",
			Description:     "Test leaderboard validation unit test",
			StartAt:         time.Now(),
			EndAt:           time.Time{},
			AggregationMode: AggregationModeMax,
			DataType:        DataTypeInt,
			Ordering:        OrderingDesc,
		}
	)

	t.Run("OK", func(t *testing.T) {
		createFunc := BuildCreateFunc(func(ctx context.Context, leaderboard Leaderboard) (string, error) {
			return expectedID, nil
		})

		id, err := createFunc(ctx, leaderboard)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, id, "wrong id returned")
	})

	t.Run("Validation Error", func(t *testing.T) {
		createFunc := BuildCreateFunc(func(ctx context.Context, leaderboard Leaderboard) (string, error) {
			return expectedID, nil
		})

		id, err := createFunc(ctx, Leaderboard{})

		assert.ErrorIs(t, err, ErrValidationError)
		assert.Empty(t, id, "id should not be returned")
	})

	t.Run("Random Error", func(t *testing.T) {
		createFunc := BuildCreateFunc(func(ctx context.Context, leaderboard Leaderboard) (string, error) {
			return "", errors.New("any error")
		})

		id, err := createFunc(ctx, leaderboard)

		assert.Error(t, err)
		assert.Empty(t, id, "id should not be returned")
	})
}
