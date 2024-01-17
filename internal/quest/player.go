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
