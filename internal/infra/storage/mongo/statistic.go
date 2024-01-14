package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/gabarcia/metagaming-api/internal/statistic"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const statisticCollectionName = "statistics"

type Statistic struct {
	CreatedAt       time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt       time.Time          `bson:"updatedAt,omitempty"`
	DeletedAt       time.Time          `bson:"deletedAt,omitempty"`
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	GameID          string             `bson:"gameId,omitempty"`
	Name            string             `bson:"name,omitempty"`
	Description     string             `bson:"description,omitempty"`
	AggregationMode string             `bson:"aggregationMode,omitempty"`
	CanOverflow     bool               `bson:"canOverflow,omitempty"`
	Goal            *float64           `bson:"goal,omitempty"`
	Landmarks       []float64          `bson:"landmarks,omitempty"`
}

func (s Statistic) toDomain() statistic.Statistic {
	return statistic.Statistic{
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		DeletedAt:       s.DeletedAt,
		ID:              s.ID.Hex(),
		GameID:          s.GameID,
		Name:            s.Name,
		Description:     s.Description,
		AggregationMode: s.AggregationMode,
		CanOverflow:     s.CanOverflow,
		Goal:            s.Goal,
		Landmarks:       s.Landmarks,
	}
}

func newStatisticFromDomain(s statistic.NewStatisticData) Statistic {
	return Statistic{
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		GameID:          s.GameID,
		Name:            s.Name,
		Description:     s.Description,
		AggregationMode: s.AggregationMode,
		CanOverflow:     s.CanOverflow,
		Goal:            s.Goal,
		Landmarks:       s.Landmarks,
	}
}

func (c connection) CreateStatistic(ctx context.Context, data statistic.NewStatisticData) (statistic.Statistic, error) {
	st := newStatisticFromDomain(data)

	cursor, err := c.client.Database(c.db).Collection(statisticCollectionName).InsertOne(ctx, st)
	if err != nil {
		return statistic.Statistic{}, err
	}

	st.ID = cursor.InsertedID.(primitive.ObjectID)

	return st.toDomain(), nil
}

func (c connection) GetStatisticByIDAndGameID(ctx context.Context, id, gameID string) (statistic.Statistic, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return statistic.Statistic{}, statistic.ErrInvalidStatisticID
	}

	cursor := c.client.Database(c.db).Collection(statisticCollectionName).FindOne(ctx, bson.M{
		"_id":       bson.M{"$eq": oid},
		"gameId":    bson.M{"$eq": gameID},
		"deletedAt": nil,
	})
	if err := cursor.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = statistic.ErrStatisticNotFound
		}

		return statistic.Statistic{}, err
	}

	var data Statistic
	if err := cursor.Decode(&data); err != nil {
		return statistic.Statistic{}, err
	}

	return data.toDomain(), nil
}

func (c connection) SoftDeleteStatistic(ctx context.Context, id, gameID string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return statistic.ErrInvalidStatisticID
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

	cursor, err := c.client.Database(c.db).Collection(statisticCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if cursor.MatchedCount == 0 {
		return statistic.ErrStatisticNotFound
	}

	return nil
}
