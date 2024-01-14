package mongo

import (
	"context"
	"time"

	"github.com/gabarcia/metagaming-api/internal/statistic"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const statisticCollectionName = "statistics"

type Statistic struct {
	CreatedAt       time.Time          `json:"createdAt,omitempty"`
	UpdatedAt       time.Time          `json:"updatedAt,omitempty"`
	DeletedAt       time.Time          `json:"deletedAt,omitempty"`
	ID              primitive.ObjectID `json:"_id,omitempty"`
	GameID          string             `json:"gameId,omitempty"`
	Name            string             `json:"name,omitempty"`
	Description     string             `json:"description,omitempty"`
	AggregationMode string             `json:"aggregationMode,omitempty"`
	CanOverflow     bool               `json:"canOverflow,omitempty"`
	Goal            *float64           `json:"goal,omitempty"`
	Landmarks       []float64          `json:"landmarks,omitempty"`
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
