package leaderboard

import "context"

type (
	// Storage function that is responsible for creating the leaderboard
	StorageCreateLeaderboardFunc func(ctx context.Context, leaderboard Leaderboard) (string, error)
)
