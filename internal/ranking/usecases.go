package ranking

import (
	"context"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"
)

type (
	// Set or update the player's rank
	UpsertPlayerRankFunc func(ctx context.Context, leaderboard leaderboard.Leaderboard, playerID string, value float64) error

	// Leaderboard ranking paginated
	RankingFunc func(ctx context.Context, leaderboard leaderboard.Leaderboard, page, limit int64) ([]Rank, error)
)
