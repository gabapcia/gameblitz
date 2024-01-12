package ranking

import "context"

type (
	// Increment player's rank value by the value provided
	StorageIncrementPlayerRankValueFunc func(ctx context.Context, leaderboardID, playerID string, value float64) error

	// Updates the player's rank value if the value provided is greather than the current one
	StorageSetMaxPlayerRankValueFunc func(ctx context.Context, leaderboardID, playerID string, value float64) error

	// Updates the player's rank value if the value provided is lower than the current one
	StorageSetMinPlayerRankValueFunc func(ctx context.Context, leaderboardID, playerID string, value float64) error

	// Get the leaderboard ranking paginated
	StorageGetRankingFunc func(ctx context.Context, leaderboardID, ordering string, page, limit int64) ([]Rank, error)
)
