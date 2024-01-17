package quest

import "context"

type (
	// Creates a quest and its tasks
	StorageCreateQuestFunc func(ctx context.Context, data NewQuestData) (Quest, error)

	// Get quest by id and game id
	StorageGetQuestFunc func(ctx context.Context, id, gameID string) (Quest, error)

	// Soft deletes a quest and its tasks
	StorageSoftDeleteQuestFunc func(ctx context.Context, questID, gameID string) error

	// Start the quest for a player
	StorageStartQuestForPlayerFunc func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error)

	// Get the player quest progression
	StorageGetPlayerQuestProgressionFunc func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error)

	// Marks the player quest task as completed and returns the tasks that were waiting for this completion
	StorageMarkPlayerQuestTaskAsCompletedFunc func(ctx context.Context, quest Quest, task Task, playerID string) ([]Task, error)
)
