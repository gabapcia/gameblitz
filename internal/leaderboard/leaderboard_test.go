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
		data := NewLeaderboardData{
			GameID:          uuid.NewString(),
			Name:            "Test Leaderboard",
			Description:     "Test leaderboard validation unit test",
			StartAt:         time.Now(),
			EndAt:           time.Time{},
			AggregationMode: AggregationModeMax,
			Ordering:        OrderingDesc,
		}

		assert.NoError(t, data.validate())
	})

	t.Run("Invalid Fields", func(t *testing.T) {
		data := NewLeaderboardData{
			GameID:          "",
			Name:            "",
			Description:     "Test leaderboard validation unit test",
			StartAt:         time.Time{},
			EndAt:           time.Time{},
			AggregationMode: "INVALID",
			Ordering:        "INVALID",
		}

		assert.ErrorIs(t, data.validate(), ErrValidationError)
		assert.ErrorIs(t, data.validate(), ErrInvalidGameID)
		assert.ErrorIs(t, data.validate(), ErrInvalidName)
		assert.ErrorIs(t, data.validate(), ErrInvalidStartDate)
		assert.ErrorIs(t, data.validate(), ErrInvalidAggregationMode)
		assert.ErrorIs(t, data.validate(), ErrInvalidOrdering)
	})

	t.Run("End Date Before Start Date", func(t *testing.T) {
		data := NewLeaderboardData{
			GameID:          uuid.NewString(),
			Name:            "Test Leaderboard",
			Description:     "Test leaderboard validation unit test",
			StartAt:         time.Now(),
			EndAt:           time.Now().Add(-24 * time.Hour),
			AggregationMode: AggregationModeMax,
			Ordering:        OrderingDesc,
		}

		assert.ErrorIs(t, data.validate(), ErrValidationError)
		assert.ErrorIs(t, data.validate(), ErrEndDateBeforeStartDate)
	})
}

func TestLeaderboardClosed(t *testing.T) {
	t.Run("Leaderboard Deleted", func(t *testing.T) {
		isClosed := Leaderboard{DeletedAt: time.Now()}.Closed()
		assert.Equal(t, true, isClosed)
	})

	t.Run("Leaderboard Not Started", func(t *testing.T) {
		isClosed := Leaderboard{StartAt: time.Now().Add(24 * time.Hour)}.Closed()
		assert.Equal(t, true, isClosed)
	})

	t.Run("Leaderboard Ended", func(t *testing.T) {
		isClosed := Leaderboard{EndAt: time.Now().Add(-24 * time.Hour)}.Closed()
		assert.Equal(t, true, isClosed)
	})

	t.Run("Leaderboard Started", func(t *testing.T) {
		isClosed := Leaderboard{StartAt: time.Now().Add(-24 * time.Hour)}.Closed()
		assert.Equal(t, false, isClosed)
	})
}

func TestBuildCreateFunc(t *testing.T) {
	var (
		ctx          = context.Background()
		expectedID   = uuid.NewString()
		expectedData = NewLeaderboardData{
			GameID:          uuid.NewString(),
			Name:            "Test Leaderboard",
			Description:     "Test create leaderboard unit test",
			StartAt:         time.Now(),
			EndAt:           time.Time{},
			AggregationMode: AggregationModeMax,
			Ordering:        OrderingDesc,
		}
	)

	t.Run("OK", func(t *testing.T) {
		createFunc := BuildCreateFunc(func(ctx context.Context, data NewLeaderboardData) (Leaderboard, error) {
			return Leaderboard{
				ID:              expectedID,
				GameID:          data.GameID,
				Name:            data.Name,
				Description:     data.Description,
				StartAt:         data.StartAt,
				EndAt:           data.EndAt,
				AggregationMode: data.AggregationMode,
				Ordering:        data.Ordering,
			}, nil
		})

		leaderboard, err := createFunc(ctx, expectedData)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, leaderboard.ID)
	})

	t.Run("Validation Error", func(t *testing.T) {
		createFunc := BuildCreateFunc(func(ctx context.Context, data NewLeaderboardData) (Leaderboard, error) {
			return Leaderboard{}, nil
		})

		leaderboard, err := createFunc(ctx, NewLeaderboardData{})

		assert.ErrorIs(t, err, ErrValidationError)
		assert.Empty(t, leaderboard.ID)
	})

	t.Run("Random Error", func(t *testing.T) {
		createFunc := BuildCreateFunc(func(ctx context.Context, data NewLeaderboardData) (Leaderboard, error) {
			return Leaderboard{}, errors.New("any error")
		})

		leaderboard, err := createFunc(ctx, expectedData)

		assert.Error(t, err)
		assert.Empty(t, leaderboard.ID)
	})
}

func TestBuildGetByIDAndGameIDFunc(t *testing.T) {
	var (
		ctx                 = context.Background()
		leaderboardID       = uuid.NewString()
		gameID              = uuid.NewString()
		expectedLeaderboard = Leaderboard{
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			ID:              leaderboardID,
			GameID:          gameID,
			Name:            "Test Leaderboard",
			Description:     "Test get leaderboard by id and game id unit test",
			StartAt:         time.Now(),
			EndAt:           time.Time{},
			AggregationMode: AggregationModeMax,
			Ordering:        OrderingDesc,
		}
	)

	t.Run("OK", func(t *testing.T) {
		getByIDAndGameIDFunc := BuildGetByIDAndGameIDFunc(func(ctx context.Context, id, gameID string) (Leaderboard, error) {
			return expectedLeaderboard, nil
		})

		leaderboard, err := getByIDAndGameIDFunc(ctx, leaderboardID, gameID)

		assert.NoError(t, err)
		assert.Equal(t, expectedLeaderboard.ID, leaderboard.ID)
	})

	t.Run("Invalid Leaderboard ID", func(t *testing.T) {
		getByIDAndGameIDFunc := BuildGetByIDAndGameIDFunc(func(ctx context.Context, id, gameID string) (Leaderboard, error) {
			return Leaderboard{}, ErrInvalidLeaderboardID
		})

		leaderboard, err := getByIDAndGameIDFunc(ctx, leaderboardID, gameID)

		assert.ErrorIs(t, err, ErrInvalidLeaderboardID)
		assert.Empty(t, leaderboard.ID)
	})

	t.Run("Not Found", func(t *testing.T) {
		getByIDAndGameIDFunc := BuildGetByIDAndGameIDFunc(func(ctx context.Context, id, gameID string) (Leaderboard, error) {
			return Leaderboard{}, ErrLeaderboardNotFound
		})

		leaderboard, err := getByIDAndGameIDFunc(ctx, leaderboardID, gameID)

		assert.ErrorIs(t, err, ErrLeaderboardNotFound)
		assert.Empty(t, leaderboard.ID)
	})

	t.Run("Random Error", func(t *testing.T) {
		getByIDAndGameIDFunc := BuildGetByIDAndGameIDFunc(func(ctx context.Context, id, gameID string) (Leaderboard, error) {
			return Leaderboard{}, errors.New("any error")
		})

		leaderboard, err := getByIDAndGameIDFunc(ctx, leaderboardID, gameID)

		assert.Error(t, err)
		assert.Empty(t, leaderboard.ID)
	})
}

func TestBuildSoftDeleteFunc(t *testing.T) {
	var (
		ctx           = context.Background()
		leaderboardID = uuid.NewString()
		gameID        = uuid.NewString()
	)

	t.Run("OK", func(t *testing.T) {
		softDeleteFunc := BuildSoftDeleteFunc(func(ctx context.Context, id, gameID string) error {
			return nil
		})

		err := softDeleteFunc(ctx, leaderboardID, gameID)

		assert.NoError(t, err)
	})

	t.Run("Invalid Leaderboard ID", func(t *testing.T) {
		softDeleteFunc := BuildSoftDeleteFunc(func(ctx context.Context, id, gameID string) error {
			return ErrInvalidLeaderboardID
		})

		err := softDeleteFunc(ctx, leaderboardID, gameID)

		assert.ErrorIs(t, err, ErrInvalidLeaderboardID)
	})

	t.Run("Not Found", func(t *testing.T) {
		softDeleteFunc := BuildSoftDeleteFunc(func(ctx context.Context, id, gameID string) error {
			return ErrLeaderboardNotFound
		})

		err := softDeleteFunc(ctx, leaderboardID, gameID)

		assert.ErrorIs(t, err, ErrLeaderboardNotFound)
	})

	t.Run("Random Error", func(t *testing.T) {
		softDeleteFunc := BuildSoftDeleteFunc(func(ctx context.Context, id, gameID string) error {
			return errors.New("any error")
		})

		err := softDeleteFunc(ctx, leaderboardID, gameID)

		assert.Error(t, err)
	})
}
