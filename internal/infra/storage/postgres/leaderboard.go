package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/gabarcia/metagaming-api/internal/infra/storage/postgres/internal/sqlc"
	"github.com/gabarcia/metagaming-api/internal/leaderboard"
)

func sqlcLeaderboardToDomain(l sqlc.Leaderboard) leaderboard.Leaderboard {
	return leaderboard.Leaderboard{
		CreatedAt:       l.CreatedAt.Time,
		UpdatedAt:       l.UpdatedAt.Time,
		DeletedAt:       l.DeletedAt.Time,
		ID:              l.ID.String(),
		GameID:          l.GameID,
		Name:            l.Name,
		Description:     l.Description,
		StartAt:         l.StartAt.Time,
		EndAt:           l.EndAt.Time,
		AggregationMode: l.AggregationMode,
		DataType:        l.DataType,
		Ordering:        l.Ordering,
	}
}

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

func (c connection) GetLeaderboardByIDAndGameID(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return leaderboard.Leaderboard{}, leaderboard.ErrInvalidLeaderboardID
	}

	data, err := c.queries.GetLeaderboardByIDAndGameID(ctx, sqlc.GetLeaderboardByIDAndGameIDParams{
		ID:     uid,
		GameID: gameID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = leaderboard.ErrLeaderboardNotFound
		}

		return leaderboard.Leaderboard{}, err
	}

	return sqlcLeaderboardToDomain(data), nil
}

func (c connection) SoftDeleteLeaderboard(ctx context.Context, id, gameID string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return leaderboard.ErrInvalidLeaderboardID
	}

	deleteCount, err := c.queries.SoftDeleteLeaderboard(ctx, sqlc.SoftDeleteLeaderboardParams{
		ID:     uid,
		GameID: gameID,
	})
	if err != nil {
		return err
	}

	if deleteCount == 0 {
		return leaderboard.ErrLeaderboardNotFound
	}

	return nil
}
