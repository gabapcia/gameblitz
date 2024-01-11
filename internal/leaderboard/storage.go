package leaderboard

import "context"

type (
	// Storage function that is responsible for creating the leaderboard
	CreateLeaderboardFunc func(ctx context.Context, leaderboard Leaderboard) (string, error)
)
