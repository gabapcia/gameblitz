package statistic

import "context"

type (
	// Create a statistic
	CreateStatisticFunc func(ctx context.Context, data NewStatisticData) (Statistic, error)

	// Get statistic by is and game id
	GetStatisticByIDAndGameID func(ctx context.Context, id, gameID string) (Statistic, error)

	// Soft delete a statistic by id and game id
	SoftDeleteStatistic func(ctx context.Context, id, gameID string) error
)
