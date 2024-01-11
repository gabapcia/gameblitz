package leaderboard

import "context"

type (
	// Storage function that is responsible for creating the leaderboard
	StorageCreateFunc func(ctx context.Context, leaderboard Leaderboard) (string, error)
)
