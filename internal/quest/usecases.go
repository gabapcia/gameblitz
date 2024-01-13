package quest

import "context"

type (
	// Creates a quest and its tasks
	CreateQuestFunc func(ctx context.Context, data NewQuestData) (Quest, error)
)
