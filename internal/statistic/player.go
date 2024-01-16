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
		Value       float64   // Landmark value
		CompletedAt time.Time // Time the player reached the landmark
	}

	PlayerProgressionUpdates struct {
		GoalJustCompleted      bool                               // Player just reached the goal?
		GoalCompletedAt        time.Time                          // Time the player reached the goal
		LandmarksJustCompleted []PlayerProgressionUpdatesLandmark // Landmarks that the player has just reached
	}
)

type (
	PlayerProgressionLandmark struct {
		Value       float64   // Landmark value
		Completed   bool      // Has the player reached the landmark?
		CompletedAt time.Time // Time the player reached the landmark
	}

	PlayerProgression struct {
		StartedAt       time.Time                   // Time the player started the progression for the given statistic
		UpdatedAt       time.Time                   // Last time the player updated it's statistic progress
		PlayerID        string                      // Player's ID
		StatisticID     string                      // Statistic ID
		CurrentValue    *float64                    // Current progression value
		GoalValue       *float64                    // Statistic's goal
		GoalCompleted   *bool                       // Has the player reached the goal?
		GoalCompletedAt time.Time                   // Time the player reached the goal
		Landmarks       []PlayerProgressionLandmark // Landmarks player progression
	}
)

func BuildUpsertPlayerProgressionFunc(
	notifierPlayerProgressionUpdates NotifierPlayerProgressionUpdates,
	storageUpdatePlayerProgressionFunc StorageUpdatePlayerProgressionFunc,
) UpsertPlayerProgressionFunc {
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

func BuildGetPlayerProgression(storageGetPlayerProgressionFunc StorageGetPlayerProgressionFunc) GetPlayerProgressionFunc {
	return func(ctx context.Context, statisticID, playerID string) (PlayerProgression, error) {
		return storageGetPlayerProgressionFunc(ctx, statisticID, playerID)
	}
}
