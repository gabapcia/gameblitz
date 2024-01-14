package statistic

import (
	"context"
	"errors"
	"time"
)

var (
	ErrPlayerStatisticNotFound = errors.New("player statistic not found")
)

type (
	PlayerProgressionUpdatesLandmark struct {
		Value       float64
		CompletedAt time.Time
	}

	PlayerProgressionUpdates struct {
		GoalJustCompleted      bool
		GoalCompletedAt        time.Time
		LandmarksJustCompleted []PlayerProgressionUpdatesLandmark
	}
)

type (
	PlayerProgressionLandmark struct {
		Value       float64
		Completed   bool
		CompletedAt time.Time
	}

	PlayerProgression struct {
		StartedAt                time.Time
		PlayerID                 string
		StatisticID              string
		StatisticAggregationMode string
		CurrentValue             *float64
		GoalValue                *float64
		GoalCompleted            *bool
		GoalCompletedAt          time.Time
		Landmarks                []PlayerProgressionLandmark
	}
)

func BuildUpdatePlayerProgressionFunc(
	notifierPlayerProgressionUpdates NotifierPlayerProgressionUpdates,
	storageUpdatePlayerProgressionFunc StorageUpdatePlayerProgressionFunc,
) UpdatePlayerProgressionFunc {
	return func(ctx context.Context, statistic Statistic, playerID string, value float64) error {
		playerProgression, playerProgressionUpdates, err := storageUpdatePlayerProgressionFunc(ctx, statistic, playerID, value)
		if err != nil {
			return err
		}

		if len(playerProgressionUpdates.LandmarksJustCompleted) > 0 || playerProgressionUpdates.GoalJustCompleted {
			return notifierPlayerProgressionUpdates(ctx, statistic, playerProgression, playerProgressionUpdates)
		}

		return nil
	}
}
