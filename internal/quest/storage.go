package quest

import "context"

type (
	// Creates a quest and its tasks
	StorageCreateQuestFunc func(ctx context.Context, data NewQuestData) (Quest, error)
)
