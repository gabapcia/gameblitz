package redis

import (
	"context"
	"fmt"

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
