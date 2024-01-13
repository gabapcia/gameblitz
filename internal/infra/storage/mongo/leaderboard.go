package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/gabarcia/metagaming-api/internal/leaderboard"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	leaderboardCollectionName = "leaderboards"
)

type Leaderboard struct {
	CreatedAt       time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt       time.Time          `bson:"updatedAt,omitempty"`
	DeletedAt       time.Time          `bson:"deletedAt,omitempty"`
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	GameID          string             `bson:"gameId,omitempty"`
	Name            string             `bson:"name,omitempty"`
	Description     string             `bson:"description,omitempty"`
	StartAt         time.Time          `bson:"startAt,omitempty"`
	EndAt           time.Time          `bson:"endAt,omitempty"`
	AggregationMode string             `bson:"aggregationMode,omitempty"`
	Ordering        string             `bson:"ordering,omitempty"`
}

func (l Leaderboard) toDomain() leaderboard.Leaderboard {
	return leaderboard.Leaderboard{
		CreatedAt:       l.CreatedAt,
		UpdatedAt:       l.UpdatedAt,
		DeletedAt:       l.DeletedAt,
		ID:              l.ID.Hex(),
		GameID:          l.GameID,
		Name:            l.Name,
		Description:     l.Description,
		StartAt:         l.StartAt,
		EndAt:           l.EndAt,
		AggregationMode: l.AggregationMode,
		Ordering:        l.Ordering,
	}
}

func newLeaderboardFromData(data leaderboard.NewLeaderboardData) Leaderboard {
	return Leaderboard{
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		GameID:          data.GameID,
		Name:            data.Name,
		Description:     data.Description,
		StartAt:         data.StartAt,
		EndAt:           data.EndAt,
		AggregationMode: data.AggregationMode,
		Ordering:        data.Ordering,
	}
}

func (c connection) CreateLeaderboard(ctx context.Context, data leaderboard.NewLeaderboardData) (leaderboard.Leaderboard, error) {
	lb := newLeaderboardFromData(data)

	cursor, err := c.client.Database(c.db).Collection(leaderboardCollectionName).InsertOne(ctx, lb)
	if err != nil {
		return leaderboard.Leaderboard{}, err
	}

	lb.ID = cursor.InsertedID.(primitive.ObjectID)

	return lb.toDomain(), nil
}

func (c connection) GetLeaderboardByIDAndGameID(ctx context.Context, id, gameID string) (leaderboard.Leaderboard, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return leaderboard.Leaderboard{}, leaderboard.ErrInvalidLeaderboardID
	}

	cursor := c.client.Database(c.db).Collection(leaderboardCollectionName).FindOne(ctx, bson.M{
		"_id":       bson.M{"$eq": oid},
		"gameId":    bson.M{"$eq": gameID},
		"deletedAt": nil,
	})
	if err = cursor.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = leaderboard.ErrLeaderboardNotFound
		}

		return leaderboard.Leaderboard{}, err
	}

	var data Leaderboard
	if err = cursor.Decode(&data); err != nil {
		return leaderboard.Leaderboard{}, err
	}

	return data.toDomain(), nil
}

func (c connection) SoftDeleteLeaderboard(ctx context.Context, id, gameID string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return leaderboard.ErrInvalidLeaderboardID
	}

	filter := bson.M{
		"_id":       bson.M{"$eq": oid},
		"gameId":    bson.M{"$eq": gameID},
		"deletedAt": nil,
	}

	update := bson.M{
		"$currentDate": bson.M{
			"deletedAt": true,
		},
	}

	cursor, err := c.client.Database(c.db).Collection(leaderboardCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if cursor.MatchedCount == 0 {
		return leaderboard.ErrLeaderboardNotFound
	}

	return nil
}
