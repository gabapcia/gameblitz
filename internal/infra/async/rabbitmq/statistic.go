package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabarcia/metagaming-api/internal/statistic"

	amqp "github.com/rabbitmq/amqp091-go"
)

const statisticExchange = "metagaming.statistic"

type (
	PlayerStatisticLandmarksUpdatesMessage struct {
		Value       float64   `json:"value"`
		CompletedAt time.Time `json:"completedAt"`
	}

	PlayerStatisticUpdatesMessage struct {
		GoalJustCompleted      bool                                     `json:"goalJustCompleted"`
		GoalCompletedAt        *time.Time                               `json:"goalCompletedAt"`
		LandmarksJustCompleted []PlayerStatisticLandmarksUpdatesMessage `json:"landmarksJustCompleted"`
	}

	PlayerStatisticLandmarkMessage struct {
		Value       float64    `json:"value"`
		Completed   bool       `json:"completed"`
		CompletedAt *time.Time `json:"completedAt"`
	}

	PlayerStatisticMessage struct {
		StartedAt       time.Time                        `json:"startedAt"`
		UpdatedAt       time.Time                        `json:"updatedAt"`
		PlayerID        string                           `json:"playerId"`
		StatisticID     string                           `json:"statisticId"`
		CurrentValue    *float64                         `json:"currentValue"`
		GoalValue       *float64                         `json:"goalValue"`
		GoalCompleted   *bool                            `json:"goalCompleted"`
		GoalCompletedAt *time.Time                       `json:"goalCompletedAt"`
		Landmarks       []PlayerStatisticLandmarkMessage `json:"landmarks"`
		LastUpdate      PlayerStatisticUpdatesMessage    `json:"lastUpdate"`
	}
)

func messageFromPlayerStatisticUpdates(progression statistic.PlayerProgression, updates statistic.PlayerProgressionUpdates) PlayerStatisticMessage {
	landmarks := make([]PlayerStatisticLandmarkMessage, len(progression.Landmarks))
	for i, landmark := range progression.Landmarks {
		var landmarkCompletedAt *time.Time
		if !landmark.CompletedAt.IsZero() {
			tmp := landmark.CompletedAt
			landmarkCompletedAt = &tmp
		}

		landmarks[i] = PlayerStatisticLandmarkMessage{
			Value:       landmark.Value,
			Completed:   landmark.Completed,
			CompletedAt: landmarkCompletedAt,
		}
	}

	landmarksUpdates := make([]PlayerStatisticLandmarksUpdatesMessage, len(updates.LandmarksJustCompleted))
	for i, landmark := range updates.LandmarksJustCompleted {
		landmarksUpdates[i] = PlayerStatisticLandmarksUpdatesMessage{
			Value:       landmark.Value,
			CompletedAt: landmark.CompletedAt,
		}
	}

	var progressionGoalCompletedAt *time.Time
	if !progression.GoalCompletedAt.IsZero() {
		progressionGoalCompletedAt = &progression.GoalCompletedAt
	}

	var updatesGoalCompletedAt *time.Time
	if !updates.GoalCompletedAt.IsZero() {
		updatesGoalCompletedAt = &updates.GoalCompletedAt
	}

	return PlayerStatisticMessage{
		StartedAt:       progression.StartedAt,
		UpdatedAt:       progression.UpdatedAt,
		PlayerID:        progression.PlayerID,
		StatisticID:     progression.StatisticID,
		CurrentValue:    progression.CurrentValue,
		GoalValue:       progression.GoalValue,
		GoalCompleted:   progression.GoalCompleted,
		GoalCompletedAt: progressionGoalCompletedAt,
		Landmarks:       landmarks,
		LastUpdate: PlayerStatisticUpdatesMessage{
			GoalJustCompleted:      updates.GoalJustCompleted,
			GoalCompletedAt:        updatesGoalCompletedAt,
			LandmarksJustCompleted: landmarksUpdates,
		},
	}
}

func buildStatisticRoutingKey(gameID, statisticID string) string {
	return fmt.Sprintf("game.%s.statistic.%s", gameID, statisticID)
}

func (p producer) ensureStatisticExchange(ctx context.Context) error {
	return p.declareExchange(ctx, statisticExchange)
}

func (p producer) PlayerStatisticProgressionUpdates(ctx context.Context, st statistic.Statistic, progression statistic.PlayerProgression, updates statistic.PlayerProgressionUpdates) error {
	var (
		routingKey = buildStatisticRoutingKey(st.GameID, st.ID)
		mandatory  = false
		immediate  = false
	)

	body, err := json.Marshal(messageFromPlayerStatisticUpdates(progression, updates))
	if err != nil {
		return err
	}

	ch, err := p.getChannel()
	if err != nil {
		return err
	}

	return ch.PublishWithContext(ctx, statisticExchange, routingKey, mandatory, immediate, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
