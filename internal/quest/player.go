package quest

import (
	"context"
	"errors"
	"time"
)

var (
	ErrPlayerAlreadyStartedTheQuest = errors.New("player already started the quest")
	ErrPlayerNotStartedTheQuest     = errors.New("player not started the quest")
	ErrPlayerQuestTaskNotRegistered = errors.New("player not registered to the task")
	ErrPlayerQuestAlreadyCompleted  = errors.New("player already concluded the quest")
)

type (
	PlayerTaskProgression struct {
		StartedAt   time.Time // Time the player started the task
		UpdatedAt   time.Time // Last time the player updated the task progression
		Task        Task      // Task config data
		CompletedAt time.Time // Time the player completed the task
	}

	PlayerQuestProgression struct {
		StartedAt        time.Time               // Time the player started the quest
		UpdatedAt        time.Time               // Last time the player updated the quest progression
		PlayerID         string                  // Player's ID
		Quest            Quest                   // Quest Config Data
		CompletedAt      time.Time               // Time the player completed the quest
		TasksProgression []PlayerTaskProgression // Tasks progression
	}
)

func (p PlayerQuestProgression) applyRuleToActiveTasks(data string) ([]string, error) {
	tasksCompleted := make([]string, 0)
	for _, taskProgression := range p.TasksProgression {
		if !taskProgression.CompletedAt.IsZero() {
			continue
		}

		pass, err := RuleApply(taskProgression.Task.Rule, data)
		if err != nil {
			return nil, err
		}

		if pass {
			tasksCompleted = append(tasksCompleted, taskProgression.Task.ID)
		}
	}

	return tasksCompleted, nil
}

func BuildStartQuestForPlayerFunc(storageStartQuestForPlayerFunc StorageStartQuestForPlayerFunc) StartQuestForPlayerFunc {
	return func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
		return storageStartQuestForPlayerFunc(ctx, quest, playerID)
	}
}

func BuildGetPlayerQuestProgression(storageGetPlayerQuestProgressionFunc StorageGetPlayerQuestProgressionFunc) GetPlayerQuestProgressionFunc {
	return func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
		return storageGetPlayerQuestProgressionFunc(ctx, quest, playerID)
	}
}

func BuildUpdatePlayerQuestProgressionFunc(
	storageGetPlayerQuestProgressionFunc StorageGetPlayerQuestProgressionFunc,
	storageUpdatePlayerQuestProgressionFunc StorageUpdatePlayerQuestProgressionFunc,
) UpdatePlayerQuestProgressionFunc {
	return func(ctx context.Context, quest Quest, playerID, taskDataToCheck string) (PlayerQuestProgression, error) {
		previousProgression, err := storageGetPlayerQuestProgressionFunc(ctx, quest, playerID)
		if err != nil {
			return PlayerQuestProgression{}, nil
		}

		if !previousProgression.CompletedAt.IsZero() {
			return PlayerQuestProgression{}, ErrPlayerQuestAlreadyCompleted
		}

		tasksCompleted, err := previousProgression.applyRuleToActiveTasks(taskDataToCheck)
		if err != nil {
			return PlayerQuestProgression{}, err
		}

		if len(tasksCompleted) == 0 {
			return previousProgression, nil
		}

		playerProgression, err := storageUpdatePlayerQuestProgressionFunc(ctx, quest, tasksCompleted, playerID)
		if err != nil {
			return PlayerQuestProgression{}, nil
		}

		return playerProgression, nil
	}
}
