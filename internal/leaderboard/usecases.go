package leaderboard

import "context"

type (
	// Create a leaderboard and return it's id
	CreateFunc func(ctx context.Context, leaderboard Leaderboard) (string, error)
)
