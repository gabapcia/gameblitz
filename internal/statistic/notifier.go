package statistic

import "context"

type (
	// Notify player progression updates
	NotifierPlayerProgressionUpdates func(ctx context.Context, statistic Statistic, updated PlayerProgressionUpdates) error
)
