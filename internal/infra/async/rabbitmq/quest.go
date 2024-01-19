package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabarcia/gameblitz/internal/quest"

	amqp "github.com/rabbitmq/amqp091-go"
)

const questExchange = "gameblitz.quest"

type (
	TaskMessage struct {
		CreatedAt             time.Time  `json:"createdAt"`
		UpdatedAt             time.Time  `json:"updatedAt"`
		DeletedAt             *time.Time `json:"deletedAt"`
		ID                    string     `json:"id"`
		Name                  string     `json:"name"`
		Description           string     `json:"description"`
		DependsOn             []string   `json:"dependsOn"`
		RequiredForCompletion bool       `json:"requiredForCompletion"`
	}

	QuestMessage struct {
		CreatedAt   time.Time     `json:"createdAt"`
		UpdatedAt   time.Time     `json:"updatedAt"`
		DeletedAt   *time.Time    `json:"deletedAt"`
		ID          string        `json:"id"`
		GameID      string        `json:"gameId"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Tasks       []TaskMessage `json:"tasks"`
	}

	PlayerTaskProgressionMessage struct {
		StartedAt   time.Time   `json:"startedAt"`
		UpdatedAt   time.Time   `json:"updatedAt"`
		Task        TaskMessage `json:"task"`
		CompletedAt *time.Time  `json:"completedAt"`
	}

	PlayerQuestProgressionMessage struct {
		StartedAt        time.Time                      `json:"startedAt"`
		UpdatedAt        time.Time                      `json:"updatedAt"`
		PlayerID         string                         `json:"playerId"`
		Quest            QuestMessage                   `json:"quest"`
		CompletedAt      *time.Time                     `json:"completedAt"`
		TasksProgression []PlayerTaskProgressionMessage `json:"tasksProgression"`
	}
)

func messageFromPlayerQuestProgression(p quest.PlayerQuestProgression) PlayerQuestProgressionMessage {
	tasks := make([]TaskMessage, len(p.Quest.Tasks))
	for i, t := range p.Quest.Tasks {
		var deletedAt *time.Time
		if !t.DeletedAt.IsZero() {
			tmp := t.DeletedAt
			deletedAt = &tmp
		}

		tasks[i] = TaskMessage{
			CreatedAt:             t.CreatedAt,
			UpdatedAt:             t.UpdatedAt,
			DeletedAt:             deletedAt,
			ID:                    t.ID,
			Name:                  t.Name,
			Description:           t.Description,
			DependsOn:             t.DependsOn,
			RequiredForCompletion: t.RequiredForCompletion,
		}
	}

	tasksProgression := make([]PlayerTaskProgressionMessage, len(p.TasksProgression))
	for i, tp := range p.TasksProgression {
		var completedAt *time.Time
		if !tp.CompletedAt.IsZero() {
			tmp := tp.CompletedAt
			completedAt = &tmp
		}

		var taskDeletedAt *time.Time
		if !tp.Task.DeletedAt.IsZero() {
			tmp := tp.Task.DeletedAt
			taskDeletedAt = &tmp
		}

		tasksProgression[i] = PlayerTaskProgressionMessage{
			StartedAt:   tp.StartedAt,
			UpdatedAt:   tp.UpdatedAt,
			CompletedAt: completedAt,
			Task: TaskMessage{
				CreatedAt:             tp.Task.CreatedAt,
				UpdatedAt:             tp.Task.UpdatedAt,
				DeletedAt:             taskDeletedAt,
				ID:                    tp.Task.ID,
				Name:                  tp.Task.Name,
				Description:           tp.Task.Description,
				DependsOn:             tp.Task.DependsOn,
				RequiredForCompletion: tp.Task.RequiredForCompletion,
			},
		}
	}

	var completedAt *time.Time
	if !p.CompletedAt.IsZero() {
		tmp := p.CompletedAt
		completedAt = &tmp
	}

	var questDeletedAt *time.Time
	if !p.Quest.DeletedAt.IsZero() {
		tmp := p.Quest.DeletedAt
		questDeletedAt = &tmp
	}

	return PlayerQuestProgressionMessage{
		StartedAt:   p.StartedAt,
		UpdatedAt:   p.UpdatedAt,
		PlayerID:    p.PlayerID,
		CompletedAt: completedAt,
		Quest: QuestMessage{
			CreatedAt:   p.Quest.CreatedAt,
			UpdatedAt:   p.Quest.UpdatedAt,
			DeletedAt:   questDeletedAt,
			ID:          p.Quest.ID,
			GameID:      p.Quest.GameID,
			Name:        p.Quest.Name,
			Description: p.Quest.Description,
			Tasks:       tasks,
		},
		TasksProgression: tasksProgression,
	}
}

func buildQuestRoutingKey(gameID, questID string) string {
	return fmt.Sprintf("game.%s.quest.%s", gameID, questID)
}

func (p producer) ensureQuestExchange(ctx context.Context) error {
	return p.declareExchange(ctx, questExchange)
}

func (p producer) PlayerQuestProgressionUpdates(ctx context.Context, progression quest.PlayerQuestProgression) error {
	var (
		routingKey = buildQuestRoutingKey(progression.Quest.GameID, progression.Quest.ID)
		mandatory  = false
		immediate  = false
	)

	body, err := json.Marshal(messageFromPlayerQuestProgression(progression))
	if err != nil {
		return err
	}

	ch, err := p.getChannel()
	if err != nil {
		return err
	}

	return ch.PublishWithContext(ctx, questExchange, routingKey, mandatory, immediate, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
