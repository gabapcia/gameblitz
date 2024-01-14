package quest

import "context"

type (
	// Creates a quest and its tasks
	StorageCreateQuestFunc func(ctx context.Context, data NewQuestData) (Quest, error)

	// Soft deletes a quest and its tasks
	StorageSoftDeleteQuestFunc func(ctx context.Context, questID, gameID string) error
)
