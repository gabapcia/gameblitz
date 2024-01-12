package ranking

import (
	"context"
	"errors"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"
)

var (
	ErrInvalidAggregationMode = errors.New("invalid aggregation mode")
	ErrInvalidOrdering        = errors.New("invalid ordering")
	ErrLeaderboardClosed      = errors.New("leaderboard closed")
	ErrInvalidPageNumber      = errors.New("invalid page number")
	ErrInvalidLimitNumber     = errors.New("invalid limit number")
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

func BuildUpsertPlayerRankFunc(
	incrementPlayerRankValueFunc StorageIncrementPlayerRankValueFunc,
	setMaxPlayerRankValueFunc StorageSetMaxPlayerRankValueFunc,
	setMinPlayerRankValueFunc StorageSetMinPlayerRankValueFunc,
) UpsertPlayerRankFunc {
	return func(ctx context.Context, lb leaderboard.Leaderboard, playerID string, value float64) error {
		if lb.Closed() {
			return ErrLeaderboardClosed
		}

		switch lb.AggregationMode {
		case leaderboard.AggregationModeInc:
			return incrementPlayerRankValueFunc(ctx, lb.ID, playerID, value)
		case leaderboard.AggregationModeMax:
			return setMaxPlayerRankValueFunc(ctx, lb.ID, playerID, value)
		case leaderboard.AggregationModeMin:
			return setMinPlayerRankValueFunc(ctx, lb.ID, playerID, value)
		default:
			return ErrInvalidAggregationMode
		}
	}
}

func BuildRankingFunc(getRankingFunc StorageGetRankingFunc) RankingFunc {
	return func(ctx context.Context, lb leaderboard.Leaderboard, page, limit int64) ([]Rank, error) {
		if page < MinPageNumber {
			return nil, ErrInvalidPageNumber
		}

		if limit < MinLimitNumber || limit > MaxLimitNumber {
			return nil, ErrInvalidLimitNumber
		}

		return getRankingFunc(ctx, lb.ID, lb.Ordering, page, limit)
	}
}
