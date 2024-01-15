package leaderboard

import "context"

type (
	// Storage function that is responsible for creating the leaderboard
	StorageCreateLeaderboardFunc func(ctx context.Context, data NewLeaderboardData) (Leaderboard, error)

	// Storage function that returns a leaderboard by it's id and game id
	StorageGetLeaderboardByIDAndGameIDFunc func(ctx context.Context, id, gameID string) (Leaderboard, error)

	// Storage function that soft delete a leaderboard
	StorageSoftDeleteLeaderboardFunc func(ctx context.Context, id, gameID string) error

	// Updates the player's rank value using the value provided
	StorageUpsertPlayerRankValueFunc func(ctx context.Context, leaderboard Leaderboard, playerID string, value float64) error

	// Get the leaderboard ranking paginated
	StorageGetRankingFunc func(ctx context.Context, leaderboardID, ordering string, page, limit int64) ([]Rank, error)
)
