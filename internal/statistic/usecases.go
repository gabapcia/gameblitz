package statistic

import "context"

type (
	// Create a statistic
	CreateFunc func(ctx context.Context, data NewStatisticData) (Statistic, error)

	// Get statistic by is and game id
	GetByIDAndGameID func(ctx context.Context, id, gameID string) (Statistic, error)

	// Soft delete a statistic by id and game id
	SoftDeleteByIDAndGameID func(ctx context.Context, id, gameID string) error

	// Update player statistic progression using the provided value
	UpdatePlayerProgressionFunc func(ctx context.Context, statistic Statistic, playerID string, value float64) error
)
