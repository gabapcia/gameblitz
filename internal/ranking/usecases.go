package ranking

import "context"

type (
	// Set or update the player's rank
	UpsertPlayerRankFunc func(ctx context.Context, leaderboardID, gameID, playerID string, value float64) error
)
