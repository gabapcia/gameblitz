package leaderboard

import "context"

type (
	// Storage function that is responsible for creating the leaderboard
	StorageCreateLeaderboardFunc func(ctx context.Context, leaderboard Leaderboard) (string, error)

	// Storage function that returns a leaderboard by it's id and game id
	StorageGetLeaderboardByIDAndGameIDFunc func(ctx context.Context, id, gameID string) (Leaderboard, error)

	// Storage function that soft delete a leaderboard
	StorageSoftDeleteLeaderboardFunc func(ctx context.Context, id, gameID string) error
)
