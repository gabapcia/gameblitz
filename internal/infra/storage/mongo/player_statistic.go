package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/gabarcia/metagaming-api/internal/statistic"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const playerStatisticCollectionName = "playersStatistics"

var (
	ErrPlayerStatisticProgressionAlreadyCreated = errors.New("player progression already created")
)

type PlayerStatisticProgressionLandmark struct {
	Value       float64   `bson:"value"`
	Completed   bool      `bson:"completed"`
	CompletedAt time.Time `bson:"completedAt,omitempty"`
}

type PlayerStatisticProgression struct {
	StartedAt                time.Time                            `bson:"startedAt,omitempty"`
	PlayerID                 string                               `bson:"playerId"`
	StatisticID              string                               `bson:"statisticId"`
	StatisticAggregationMode string                               `bson:"statisticAggregationMode"`
	CurrentValue             *float64                             `bson:"currentValue"`
	GoalValue                *float64                             `bson:"goalValue,omitempty"`
	GoalCompleted            *bool                                `bson:"goalCompleted,omitempty"`
	GoalCompletedAt          time.Time                            `bson:"goalCompletedAt,omitempty"`
	Landmarks                []PlayerStatisticProgressionLandmark `bson:"landmarks,omitempty"`

	PreviousData *PlayerStatisticProgression `bson:"_previousData,omitempty"`
}

func (p PlayerStatisticProgression) goalJustCompleted() bool {
	if p.PreviousData == nil {
		return false
	}

	var (
		prev = p.PreviousData.GoalCompleted
		curr = p.GoalCompleted
	)

	if prev == nil || curr == nil {
		return false
	}

	return !*prev && *curr
}

func (p PlayerStatisticProgression) landmarksJustCompleted() []PlayerStatisticProgressionLandmark {
	if p.PreviousData == nil {
		return nil
	}

	var (
		prevLandmarks = p.PreviousData.Landmarks
		currLandmarks = p.Landmarks
	)

	if prevLandmarks == nil || currLandmarks == nil {
		return nil
	}

	justCompleted := make([]PlayerStatisticProgressionLandmark, 0)
	for _, curr := range currLandmarks {
		for _, prev := range prevLandmarks {
			if curr.Value == prev.Value && (!prev.Completed && curr.Completed) {
				justCompleted = append(justCompleted, curr)
				break
			}
		}
	}

	return justCompleted
}

func (p PlayerStatisticProgression) toDomainUpdates() statistic.PlayerProgressionUpdates {
	var (
		landmarksJustCompletedData = p.landmarksJustCompleted()
		landmarksJustCompleted     = make([]statistic.PlayerProgressionUpdatesLandmark, len(landmarksJustCompletedData))
	)
	for i, landmark := range landmarksJustCompletedData {
		landmarksJustCompleted[i] = statistic.PlayerProgressionUpdatesLandmark{
			Value:       landmark.Value,
			CompletedAt: landmark.CompletedAt,
		}
	}

	return statistic.PlayerProgressionUpdates{
		GoalJustCompleted:      p.goalJustCompleted(),
		GoalCompletedAt:        p.GoalCompletedAt,
		LandmarksJustCompleted: landmarksJustCompleted,
	}
}

func (p PlayerStatisticProgression) toDomain() statistic.PlayerProgression {
	landmarks := make([]statistic.PlayerProgressionLandmark, len(p.Landmarks))
	for i, landmark := range p.Landmarks {
		landmarks[i] = statistic.PlayerProgressionLandmark{
			Value:       landmark.Value,
			Completed:   landmark.Completed,
			CompletedAt: landmark.CompletedAt,
		}
	}

	return statistic.PlayerProgression{
		StartedAt:                p.StartedAt,
		PlayerID:                 p.PlayerID,
		StatisticID:              p.StatisticID,
		StatisticAggregationMode: p.StatisticAggregationMode,
		CurrentValue:             p.CurrentValue,
		GoalValue:                p.GoalValue,
		GoalCompleted:            p.GoalCompleted,
		GoalCompletedAt:          p.GoalCompletedAt,
		Landmarks:                landmarks,
	}
}

func (c connection) ensurePlayerStatisticIndexes(ctx context.Context) error {
	_, err := c.client.Database(c.db).Collection(playerStatisticCollectionName).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "playerId", Value: 1},
				{Key: "statisticId", Value: 1},
			},
			Options: options.Index().SetName("playerId_1_statisticId_1").SetUnique(true),
		},
	})

	return err
}

func (c connection) createPlayerStatisticProgression(ctx context.Context, st statistic.Statistic, playerID string) error {
	var (
		goalValue     = st.Goal
		goalCompleted *bool
	)
	if goalValue != nil {
		tmpCompleted := false
		goalCompleted = &tmpCompleted
	}

	landmarks := make([]PlayerStatisticProgressionLandmark, len(st.Landmarks))
	for i, landmark := range st.Landmarks {
		landmarks[i] = PlayerStatisticProgressionLandmark{Value: landmark}
	}

	data := PlayerStatisticProgression{
		PlayerID:                 playerID,
		StatisticID:              st.ID,
		StatisticAggregationMode: st.AggregationMode,
		CurrentValue:             st.InitialValue,
		GoalValue:                goalValue,
		GoalCompleted:            goalCompleted,
		Landmarks:                landmarks,
	}

	if _, err := c.client.Database(c.db).Collection(playerStatisticCollectionName).InsertOne(ctx, data); err != nil {
		return err
	}

	return nil
}

func (c connection) getPlayerStatisticProgression(ctx context.Context, statisticID, playerID string) (PlayerStatisticProgression, error) {
	cursor := c.client.Database(c.db).Collection(playerStatisticCollectionName).FindOne(ctx, bson.M{
		"statisticId": bson.M{"$eq": statisticID},
		"playerId":    bson.M{"$eq": playerID},
	})
	if err := cursor.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = statistic.ErrPlayerStatisticNotFound
		}

		return PlayerStatisticProgression{}, err
	}

	var data PlayerStatisticProgression
	return data, cursor.Decode(&data)
}

func (c connection) updatePlayerStatisticProgression(ctx context.Context, statisticID, playerID string, value float64) (PlayerStatisticProgression, error) {
	data, err := c.getPlayerStatisticProgression(ctx, statisticID, playerID)
	if err != nil {
		return PlayerStatisticProgression{}, err
	}

	var (
		comparisonOp        = ""
		aggregationOp       = ""
		defaultCurrentValue = value
	)
	switch data.StatisticAggregationMode {
	case statistic.AggregationModeSum:
		comparisonOp = "$gte"
		aggregationOp = "$add"
		defaultCurrentValue = 0
	case statistic.AggregationModeMax:
		comparisonOp = "$gte"
		aggregationOp = "$max"
	case statistic.AggregationModeSub:
		comparisonOp = "$lte"
		aggregationOp = "$subtract"
		defaultCurrentValue = 0
	case statistic.AggregationModeMin:
		comparisonOp = "$lte"
		aggregationOp = "$min"
	default:
		return PlayerStatisticProgression{}, statistic.ErrInvalidAggregationMode
	}

	var currentValueAgg = bson.M{aggregationOp: bson.A{bson.M{"$ifNull": bson.A{"$currentValue", defaultCurrentValue}}, value}}

	filter := bson.M{
		"playerId":    bson.M{"$eq": playerID},
		"statisticId": bson.M{"$eq": statisticID},
	}

	update := bson.A{
		bson.M{"$set": bson.M{"_previousData": "$$ROOT"}},
		bson.M{"$set": bson.M{
			"startedAt": bson.M{"$cond": bson.M{
				"if": bson.M{"$eq": bson.A{
					bson.M{"$ifNull": bson.A{"$startedAt", "NULL"}},
					"NULL",
				}},
				"then": time.Now().UTC(),
				"else": "$startedAt",
			}},
			"currentValue": currentValueAgg,
			"landmarks": bson.M{"$map": bson.M{
				"input": "$landmarks",
				"as":    "landmark",
				"in": bson.M{"$cond": bson.M{
					"if": bson.M{"$and": bson.A{
						bson.M{"$eq": bson.A{"$$landmark.completed", false}},
						bson.M{comparisonOp: bson.A{currentValueAgg, "$$landmark.value"}},
					}},
					"then": bson.M{
						"$mergeObjects": bson.A{"$$landmark", bson.M{
							"completed":   true,
							"completedAt": time.Now().UTC(),
						}},
					},
					"else": "$$landmark",
				}},
			}},
			"goalCompleted": bson.M{"$cond": bson.M{
				"if": bson.M{"$eq": bson.A{
					bson.M{"$ifNull": bson.A{"$goalCompleted", "NULL"}},
					"NULL",
				}},
				"then": nil,
				"else": bson.M{"$cond": bson.M{
					"if": bson.M{"$and": bson.A{
						bson.M{"$eq": bson.A{"$goalCompleted", false}},
						bson.M{comparisonOp: bson.A{currentValueAgg, "$goalValue"}},
					}},
					"then": true,
					"else": "$goalCompleted",
				}},
			}},
			"goalCompletedAt": bson.M{"$cond": bson.M{
				"if": bson.M{"$eq": bson.A{
					bson.M{"$ifNull": bson.A{"$goalCompleted", "NULL"}},
					"NULL",
				}},
				"then": nil,
				"else": bson.M{"$cond": bson.M{
					"if": bson.M{"$and": bson.A{
						bson.M{"$eq": bson.A{"$goalCompleted", false}},
						bson.M{comparisonOp: bson.A{currentValueAgg, "$goalValue"}},
					}},
					"then": time.Now().UTC(),
					"else": "$goalCompletedAt",
				}},
			}},
		}},
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	cursor := c.client.Database(c.db).Collection(playerStatisticCollectionName).FindOneAndUpdate(ctx, filter, update, opts)
	if err := cursor.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = statistic.ErrPlayerStatisticNotFound
		}

		return PlayerStatisticProgression{}, err
	}

	data = PlayerStatisticProgression{}
	return data, cursor.Decode(&data)
}

func (c connection) upsertPlayerStatisticProgression(ctx context.Context, st statistic.Statistic, playerID string, value float64) (PlayerStatisticProgression, error) {
	progression, err := c.updatePlayerStatisticProgression(ctx, st.ID, playerID, value)
	if err != nil {
		if errors.Is(err, statistic.ErrPlayerStatisticNotFound) {
			if err := c.createPlayerStatisticProgression(ctx, st, playerID); err != nil && !errors.Is(err, ErrPlayerStatisticProgressionAlreadyCreated) {
				return PlayerStatisticProgression{}, err
			}

			return c.updatePlayerStatisticProgression(ctx, st.ID, playerID, value)
		}

		return PlayerStatisticProgression{}, err
	}

	return progression, nil
}

func (c connection) UpdatePlayerStatisticProgression(ctx context.Context, st statistic.Statistic, playerID string, value float64) (statistic.PlayerProgression, statistic.PlayerProgressionUpdates, error) {
	playerProgression, err := c.upsertPlayerStatisticProgression(ctx, st, playerID, value)
	if err != nil {
		return statistic.PlayerProgression{}, statistic.PlayerProgressionUpdates{}, err
	}

	return playerProgression.toDomain(), playerProgression.toDomainUpdates(), nil
}
