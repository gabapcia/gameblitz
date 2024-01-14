package statistic

import "context"

type (
	// Notify player progression updates
	NotifierPlayerProgressionUpdates func(ctx context.Context, playerID string, statistic Statistic, updated PlayerProgressionUpdates) error
)
