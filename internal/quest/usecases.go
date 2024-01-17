package quest

import "context"

type (
	// Creates a quest and its tasks
	CreateQuestFunc func(ctx context.Context, data NewQuestData) (Quest, error)

	// Get quest by id and game id
	GetQuestByIDAndGameIDFunc func(ctx context.Context, id, gameID string) (Quest, error)

	// Soft deletes a quest and its tasks
	SoftDeleteQuestFunc func(ctx context.Context, questID, gameID string) error

	// Start the quest for a player
	StartQuestForPlayerFunc func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error)

	// Get the player quest progression
	GetPlayerQuestProgressionFunc func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error)

	// Apply `taskDataToCheck` to all active tasks, check if it meets your conditions and update the completion of tasks that do.
	// When all the required tasks are marked as completed, the quest will also be automatically marked as completed
	UpdatePlayerQuestProgressionFunc func(ctx context.Context, quest Quest, playerID, taskDataToCheck string) (PlayerQuestProgression, error)
)
