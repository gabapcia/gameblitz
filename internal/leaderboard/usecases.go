package leaderboard

import "context"

type (
	// Create a leaderboard and return it's id
	CreateFunc func(ctx context.Context, data NewLeaderboardData) (Leaderboard, error)

	// Get a leaderboard by it's id and game id
	GetByIDAndGameIDFunc func(ctx context.Context, id, gameID string) (Leaderboard, error)

	// Soft Delete a leaderboard
	SoftDeleteFunc func(ctx context.Context, id, gameID string) error

	// Set or update the player's rank
	UpsertPlayerRankFunc func(ctx context.Context, leaderboard Leaderboard, playerID string, value float64) error

	// Leaderboard ranking paginated
	RankingFunc func(ctx context.Context, leaderboard Leaderboard, page, limit int64) ([]Rank, error)
)
