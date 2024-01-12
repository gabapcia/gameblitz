package redis

import (
	"context"
	"fmt"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"
	"github.com/gabarcia/metagaming-api/internal/ranking"

	"github.com/redis/go-redis/v9"
)

func buildRankingKey(leaderboardID string) string {
	return fmt.Sprintf("ranking:%s", leaderboardID)
}

func (c connection) IncrementPlayerRankValue(ctx context.Context, leaderboardID, playerID string, value float64) error {
	cursor := c.rdb.ZIncrBy(ctx, buildRankingKey(leaderboardID), value, playerID)
	return cursor.Err()
}

func (c connection) SetMaxPlayerRankValue(ctx context.Context, leaderboardID, playerID string, value float64) error {
	cursor := c.rdb.ZAddGT(ctx, buildRankingKey(leaderboardID), redis.Z{Score: value, Member: playerID})
	return cursor.Err()
}

func (c connection) SetMinPlayerRankValue(ctx context.Context, leaderboardID, playerID string, value float64) error {
	cursor := c.rdb.ZAddLT(ctx, buildRankingKey(leaderboardID), redis.Z{Score: value, Member: playerID})
	return cursor.Err()
}

func (c connection) GetRanking(ctx context.Context, leaderboardID, ordering string, page, limit int64) ([]ranking.Rank, error) {
	var cursor *redis.ZSliceCmd
	switch ordering {
	case leaderboard.OrderingAsc:
		cursor = c.rdb.ZRangeWithScores(ctx, buildRankingKey(leaderboardID), page*limit, limit-1)
	case leaderboard.OrderingDesc:
		cursor = c.rdb.ZRevRangeWithScores(ctx, buildRankingKey(leaderboardID), page*limit, limit-1)
	default:
		return nil, ranking.ErrInvalidOrdering
	}

	data, err := cursor.Result()
	if err != nil {
		return nil, err
	}

	rankingFiltered := make([]ranking.Rank, len(data))
	for i, d := range data {
		rankingFiltered[i] = ranking.Rank{
			LeaderboardID: leaderboardID,
			PlayerID:      d.Member.(string),
			Position:      int64(i),
			Value:         d.Score,
		}
	}

	return rankingFiltered, nil
}
