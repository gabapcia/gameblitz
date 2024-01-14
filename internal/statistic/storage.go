package statistic

import "context"

type (
	// Create a statistic
	StorageCreateStatisticFunc func(ctx context.Context, data NewStatisticData) (Statistic, error)
)
