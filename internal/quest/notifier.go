package quest

import "context"

type (
	// Notify player progression updates
	NotifierPlayerProgressionUpdates func(ctx context.Context, progression PlayerQuestProgression) error
)
