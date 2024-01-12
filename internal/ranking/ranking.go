package ranking

import (
	"context"
	"errors"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"
)

var (
	ErrInvalidAggregationMode = errors.New("invalid aggregation mode")
	ErrLeaderboardClosed      = errors.New("leaderboard closed")
)

type Rank struct {
	LeaderboardID string
	PlayerID      string
	Position      int64
	Value         float64
}

func BuildUpsertPlayerRankFunc(
	getLeaderboardByIDAndGameIDFunc leaderboard.GetByIDAndGameIDFunc,
	incrementPlayerRankValueFunc StorageIncrementPlayerRankValueFunc,
	setMaxPlayerRankValueFunc StorageSetMaxPlayerRankValueFunc,
	setMinPlayerRankValueFunc StorageSetMinPlayerRankValueFunc,
) UpsertPlayerRankFunc {
	return func(ctx context.Context, leaderboardID, gameID, playerID string, value float64) error {
		lb, err := getLeaderboardByIDAndGameIDFunc(ctx, leaderboardID, gameID)
		if err != nil {
			return err
		}

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
