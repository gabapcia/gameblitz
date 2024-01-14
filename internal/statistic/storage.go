package statistic

import "context"

type (
	// Create a statistic
	StorageCreateStatisticFunc func(ctx context.Context, data NewStatisticData) (Statistic, error)

	// Get statistic by is and game id
	StorageGetStatisticByIDAndGameID func(ctx context.Context, id, gameID string) (Statistic, error)

	// Soft delete a statistic by id and game id
	StorageSoftDeleteStatistic func(ctx context.Context, id, gameID string) error

	// Updates the player statistic progression using the provided value
	StorageUpdatePlayerProgressionFunc func(ctx context.Context, statistic Statistic, playerID string, value float64) (PlayerProgression, PlayerProgressionUpdates, error)
)
