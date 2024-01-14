package statistic

import "context"

type PlayerProgressionUpdates struct {
	PlayerID           string    // Player ID
	CurrentValue       float64   // The current statistic value for the user
	GoalComplted       bool      // Has the goal just been achieved?
	LandmarksCompleted []float64 // Landmarks that just have been reached
}

func BuildUpdatePlayerProgressionFunc(
	notifierPlayerProgressionUpdates NotifierPlayerProgressionUpdates,
	storageIncreasePlayerProgressionFunc StorageIncreasePlayerProgressionFunc,
	storageSetMaxPlayerProgressionFunc StorageSetMaxPlayerProgressionFunc,
	storageSetMinPlayerProgressionFunc StorageSetMinPlayerProgressionFunc,
) UpdatePlayerProgressionFunc {
	return func(ctx context.Context, statistic Statistic, playerID string, value float64) error {
		var (
			playerProgressionUpdates PlayerProgressionUpdates
			err                      error
		)

		switch statistic.AggregationMode {
		case AggregationModeMax:
			playerProgressionUpdates, err = storageSetMaxPlayerProgressionFunc(ctx, statistic.ID, statistic.GameID, playerID, value)
		case AggregationModeMin:
			playerProgressionUpdates, err = storageSetMinPlayerProgressionFunc(ctx, statistic.ID, statistic.GameID, playerID, value)
		case AggregationModeInc:
			playerProgressionUpdates, err = storageIncreasePlayerProgressionFunc(ctx, statistic.ID, statistic.GameID, playerID, value)
		default:
			return ErrInvalidAggregationMode
		}

		if err != nil {
			return err
		}

		if len(playerProgressionUpdates.LandmarksCompleted) > 0 || playerProgressionUpdates.GoalComplted {
			return notifierPlayerProgressionUpdates(ctx, statistic, playerProgressionUpdates)
		}

		return nil
	}
}
