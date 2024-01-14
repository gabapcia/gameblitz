package statistic

import "context"

type (
	// Create a statistic
	CreateStatisticFunc func(ctx context.Context, data NewStatisticData) (Statistic, error)
)
