package postgres

import (
	"context"
	"errors"

	"github.com/gabarcia/metagaming-api/internal/infra/storage/postgres/internal/sqlc"
	"github.com/gabarcia/metagaming-api/internal/leaderboard"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
		Ordering:        l.Ordering,
	}
}

func (c connection) CreateLeaderboard(ctx context.Context, data leaderboard.NewLeaderboardData) (leaderboard.Leaderboard, error) {
	leaderboard, err := c.queries.CreateLeaderboard(ctx, sqlc.CreateLeaderboardParams{
		GameID:          data.GameID,
		Name:            data.Name,
		Description:     data.Description,
		StartAt:         pgtype.Timestamptz{Time: data.StartAt, Valid: true},
		EndAt:           pgtype.Timestamptz{Time: data.EndAt, Valid: !data.EndAt.IsZero()},
		AggregationMode: data.AggregationMode,
		Ordering:        data.Ordering,
	})

	return sqlcLeaderboardToDomain(leaderboard), err
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
