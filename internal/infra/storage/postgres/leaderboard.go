package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/gabarcia/metagaming-api/internal/infra/storage/postgres/internal/sqlc"
	"github.com/gabarcia/metagaming-api/internal/leaderboard"
)

func (c connection) CreateLeaderboard(ctx context.Context, leaderboard leaderboard.Leaderboard) (string, error) {
	leaderboardID, err := c.queries.CreateLeaderboard(ctx, sqlc.CreateLeaderboardParams{
		GameID:          leaderboard.GameID,
		Name:            leaderboard.Name,
		Description:     leaderboard.Description,
		StartAt:         pgtype.Timestamptz{Time: leaderboard.StartAt, Valid: true},
		EndAt:           pgtype.Timestamptz{Time: leaderboard.EndAt, Valid: !leaderboard.EndAt.IsZero()},
		AggregationMode: leaderboard.AggregationMode,
		DataType:        leaderboard.DataType,
		Ordering:        leaderboard.Ordering,
	})

	return leaderboardID.String(), err
}
