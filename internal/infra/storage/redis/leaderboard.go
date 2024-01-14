package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"
	"github.com/google/uuid"
)

type Leaderboard struct {
	CreatedAt       time.Time  `redis:"createdAt,omitempty"`
	UpdatedAt       time.Time  `redis:"updatedAt,omitempty"`
	DeletedAt       *time.Time `redis:"deletedAt,omitempty"`
	ID              string     `redis:"id,omitempty"`
	GameID          string     `redis:"gameId,omitempty"`
	Name            string     `redis:"name,omitempty"`
	Description     string     `redis:"description,omitempty"`
	StartAt         time.Time  `redis:"startAt,omitempty"`
	EndAt           *time.Time `redis:"endAt,omitempty"`
	AggregationMode string     `redis:"aggregationMode,omitempty"`
	Ordering        string     `redis:"ordering,omitempty"`
}

func (l Leaderboard) toDomain() leaderboard.Leaderboard {
	var deletedAt time.Time
	if l.DeletedAt != nil {
		deletedAt = *l.DeletedAt
	}

	var endAt time.Time
	if l.EndAt != nil {
		endAt = *l.EndAt
	}

	return leaderboard.Leaderboard{
		CreatedAt:       l.CreatedAt,
		UpdatedAt:       l.UpdatedAt,
		DeletedAt:       deletedAt,
		ID:              l.ID,
		GameID:          l.GameID,
		Name:            l.Name,
		Description:     l.Description,
		StartAt:         l.StartAt,
		EndAt:           endAt,
		AggregationMode: l.AggregationMode,
		Ordering:        l.Ordering,
	}
}

func newLeaderboardFromData(data leaderboard.NewLeaderboardData) Leaderboard {
	var endAt *time.Time
	if !data.EndAt.IsZero() {
		endAt = &data.EndAt
	}

	return Leaderboard{
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		ID:              uuid.NewString(),
		GameID:          data.GameID,
		Name:            data.Name,
		Description:     data.Description,
		StartAt:         data.StartAt,
		EndAt:           endAt,
		AggregationMode: data.AggregationMode,
		Ordering:        data.Ordering,
	}
}

func buildLeaderboardKey(leaderboardID string) string {
	return fmt.Sprintf("leaderboard:%s", leaderboardID)
}

func (c connection) CreateLeaderboard(ctx context.Context, data leaderboard.NewLeaderboardData) (leaderboard.Leaderboard, error) {
	lb := newLeaderboardFromData(data)

	if err := c.rdb.HSet(ctx, buildLeaderboardKey(lb.ID), lb).Err(); err != nil {
		return leaderboard.Leaderboard{}, err
	}

	return lb.toDomain(), nil
}

func (c connection) GetLeaderboardByIDAndGameID(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
	cursor := c.rdb.HGetAll(ctx, buildLeaderboardKey(id))
	if err := cursor.Err(); err != nil {
		return leaderboard.Leaderboard{}, nil
	}

	var lb Leaderboard
	if err := cursor.Scan(&lb); err != nil {
		return leaderboard.Leaderboard{}, err
	}

	if lb.ID == "" || lb.DeletedAt != nil || lb.GameID != gameID {
		return leaderboard.Leaderboard{}, leaderboard.ErrLeaderboardNotFound
	}

	return lb.toDomain(), nil
}

func (c connection) SoftDeleteLeaderboard(ctx context.Context, id, gameID string) error {
	if _, err := c.GetLeaderboardByIDAndGameID(ctx, id, gameID); err != nil {
		return err
	}

	return c.rdb.HSetNX(ctx, buildLeaderboardKey(id), "deletedAt", time.Now().UTC()).Err()
}
