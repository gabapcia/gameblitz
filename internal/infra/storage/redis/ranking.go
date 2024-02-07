package redis

import (
	"context"
	"fmt"

	"github.com/gabapcia/gameblitz/internal/leaderboard"

	"github.com/redis/go-redis/v9"
)

func buildRankingKey(leaderboardID string) string {
	return fmt.Sprintf("leaderboard:%s:ranking", leaderboardID)
}

func (c connection) incrementPlayerRankValue(ctx context.Context, leaderboardID, playerID string, value float64) error {
	cursor := c.rdb.ZIncrBy(ctx, buildRankingKey(leaderboardID), value, playerID)
	return cursor.Err()
}

func (c connection) setMaxPlayerRankValue(ctx context.Context, leaderboardID, playerID string, value float64) error {
	cursor := c.rdb.ZAddGT(ctx, buildRankingKey(leaderboardID), redis.Z{Score: value, Member: playerID})
	return cursor.Err()
}

func (c connection) setMinPlayerRankValue(ctx context.Context, leaderboardID, playerID string, value float64) error {
	cursor := c.rdb.ZAddLT(ctx, buildRankingKey(leaderboardID), redis.Z{Score: value, Member: playerID})
	return cursor.Err()
}

func (c connection) UpsertPlayerRankValue(ctx context.Context, lb leaderboard.Leaderboard, playerID string, value float64) error {
	switch lb.AggregationMode {
	case leaderboard.AggregationModeInc:
		return c.incrementPlayerRankValue(ctx, lb.ID, playerID, value)
	case leaderboard.AggregationModeMax:
		return c.setMaxPlayerRankValue(ctx, lb.ID, playerID, value)
	case leaderboard.AggregationModeMin:
		return c.setMinPlayerRankValue(ctx, lb.ID, playerID, value)
	default:
		return leaderboard.ErrInvalidAggregationMode
	}
}

func (c connection) GetRanking(ctx context.Context, leaderboardID, ordering string, page, limit int64) ([]leaderboard.Rank, error) {
	var cursor *redis.ZSliceCmd
	switch ordering {
	case leaderboard.OrderingAsc:
		cursor = c.rdb.ZRangeWithScores(ctx, buildRankingKey(leaderboardID), page*limit, limit-1)
	case leaderboard.OrderingDesc:
		cursor = c.rdb.ZRevRangeWithScores(ctx, buildRankingKey(leaderboardID), page*limit, limit-1)
	default:
		return nil, leaderboard.ErrInvalidOrdering
	}

	data, err := cursor.Result()
	if err != nil {
		return nil, err
	}

	rankingFiltered := make([]leaderboard.Rank, len(data))
	for i, d := range data {
		rankingFiltered[i] = leaderboard.Rank{
			LeaderboardID: leaderboardID,
			PlayerID:      d.Member.(string),
			Position:      int64(i),
			Value:         d.Score,
		}
	}

	return rankingFiltered, nil
}
