package leaderboard

import (
	"context"
	"errors"
)

var (
	ErrLeaderboardClosed  = errors.New("leaderboard closed")
	ErrInvalidPageNumber  = errors.New("invalid page number")
	ErrInvalidLimitNumber = errors.New("invalid limit number")
)

const (
	MaxLimitNumber = 500
	MinLimitNumber = 1
	MinPageNumber  = 0
)

type Rank struct {
	LeaderboardID string
	PlayerID      string
	Position      int64
	Value         float64
}

func BuildUpsertPlayerRankFunc(upsertPlayerRankValueFunc StorageUpsertPlayerRankValueFunc) UpsertPlayerRankFunc {
	return func(ctx context.Context, lb Leaderboard, playerID string, value float64) error {
		if lb.Closed() {
			return ErrLeaderboardClosed
		}

		return upsertPlayerRankValueFunc(ctx, lb, playerID, value)
	}
}

func BuildRankingFunc(getRankingFunc StorageGetRankingFunc) RankingFunc {
	return func(ctx context.Context, lb Leaderboard, page, limit int64) ([]Rank, error) {
		if page < MinPageNumber {
			return nil, ErrInvalidPageNumber
		}

		if limit < MinLimitNumber || limit > MaxLimitNumber {
			return nil, ErrInvalidLimitNumber
		}

		return getRankingFunc(ctx, lb.ID, lb.Ordering, page, limit)
	}
}
