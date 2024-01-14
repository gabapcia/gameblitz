package quest

import "context"

type (
	// Creates a quest and its tasks
	CreateQuestFunc func(ctx context.Context, data NewQuestData) (Quest, error)

	// Get quest by id and game id
	GetQuestByIDAndGameIDFunc func(ctx context.Context, id, gameID string) (Quest, error)

	// Soft deletes a quest and its tasks
	SoftDeleteQuestFunc func(ctx context.Context, questID, gameID string) error
)
