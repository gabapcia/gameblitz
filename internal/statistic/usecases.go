package statistic

import "context"

type (
	// Create a statistic
	CreateFunc func(ctx context.Context, data NewStatisticData) (Statistic, error)

	// Get statistic by is and game id
	GetByIDAndGameIDFunc func(ctx context.Context, id, gameID string) (Statistic, error)

	// Soft delete a statistic by id and game id
	SoftDeleteByIDAndGameIDFunc func(ctx context.Context, id, gameID string) error

	// Update player statistic progression using the provided value
	UpsertPlayerProgressionFunc func(ctx context.Context, statistic Statistic, playerID string, value float64) error

	// Get player progression by statistic id and player id
	GetPlayerProgressionFunc func(ctx context.Context, statisticID, playerID string) (PlayerProgression, error)
)
