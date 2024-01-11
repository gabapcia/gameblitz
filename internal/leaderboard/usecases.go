package leaderboard

import "context"

type (
	// Create a leaderboard and return it's id
	CreateFunc func(ctx context.Context, data NewLeaderboardData) (Leaderboard, error)

	// Get a leaderboard by it's id and game id
	GetByIDAndGameIDFunc func(ctx context.Context, id, gameID string) (Leaderboard, error)

	// Soft Delete a leaderboard
	SoftDeleteFunc func(ctx context.Context, id, gameID string) error
)
