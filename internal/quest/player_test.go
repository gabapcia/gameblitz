package quest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlayerProgression_ApplyRuleToActiveTasks(t *testing.T) {
	t.Run("With All Tasks Completed", func(t *testing.T) {
		progression := PlayerQuestProgression{TasksProgression: []PlayerTaskProgression{
			{CompletedAt: time.Now()},
			{CompletedAt: time.Now()},
			{CompletedAt: time.Now()},
			{CompletedAt: time.Now()},
		}}

		tasksCompleted, err := progression.applyRuleToActiveTasks("")
		assert.NoError(t, err)
		assert.Empty(t, tasksCompleted)
	})

	t.Run("With Just One Pending Task", func(t *testing.T) {
		pendingTaskID := uuid.NewString()
		progression := PlayerQuestProgression{TasksProgression: []PlayerTaskProgression{
			{Task: Task{ID: uuid.NewString()}, CompletedAt: time.Now()},
			{Task: Task{ID: uuid.NewString()}, CompletedAt: time.Now()},
			{Task: Task{ID: pendingTaskID, Rule: `{"==": [{"var": "fields.bool"}, true]}`}},
			{Task: Task{ID: uuid.NewString()}, CompletedAt: time.Now()},
		}}

		tasksCompleted, err := progression.applyRuleToActiveTasks(`{"fields": {"bool": true}}`)
		assert.NoError(t, err)
		if assert.Len(t, tasksCompleted, 1) {
			assert.Equal(t, pendingTaskID, tasksCompleted[0])
		}
	})

	t.Run("With Rule Apply Error", func(t *testing.T) {
		progression := PlayerQuestProgression{TasksProgression: []PlayerTaskProgression{
			{Task: Task{ID: uuid.NewString(), Rule: `{"==": [{"var": "fields.bool"}, true]}`}},
		}}

		tasksCompleted, err := progression.applyRuleToActiveTasks(`{`)
		assert.ErrorIs(t, err, ErrBrokenRuleData)
		assert.Empty(t, tasksCompleted)
	})
}

func TestBuildStartQuestForPlayerFunc(t *testing.T) {
	var (
		ctx = context.Background()

		playerID = uuid.NewString()
		quest    = Quest{ID: uuid.NewString()}
	)

	t.Run("OK", func(t *testing.T) {
		startQuestForPlayerFunc := BuildStartQuestForPlayerFunc(func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
			return PlayerQuestProgression{Quest: quest, PlayerID: playerID}, nil
		})

		playerProgression, err := startQuestForPlayerFunc(ctx, quest, playerID)
		assert.NoError(t, err)

		assert.Equal(t, playerID, playerProgression.PlayerID)
		assert.Equal(t, quest.ID, playerProgression.Quest.ID)
	})

	t.Run("Quest Already Started For Player", func(t *testing.T) {
		startQuestForPlayerFunc := BuildStartQuestForPlayerFunc(func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
			return PlayerQuestProgression{}, ErrPlayerAlreadyStartedTheQuest
		})

		_, err := startQuestForPlayerFunc(ctx, quest, playerID)
		assert.ErrorIs(t, err, ErrPlayerAlreadyStartedTheQuest)
	})

	t.Run("Quest Not Found", func(t *testing.T) {
		startQuestForPlayerFunc := BuildStartQuestForPlayerFunc(func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
			return PlayerQuestProgression{}, ErrQuestNotFound
		})

		_, err := startQuestForPlayerFunc(ctx, quest, playerID)
		assert.ErrorIs(t, err, ErrQuestNotFound)
	})

	t.Run("Random Error", func(t *testing.T) {
		startQuestForPlayerFunc := BuildStartQuestForPlayerFunc(func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
			return PlayerQuestProgression{}, errors.New("ant error")
		})

		_, err := startQuestForPlayerFunc(ctx, quest, playerID)
		assert.Error(t, err)
	})
}

func TestBuildGetPlayerQuestProgression(t *testing.T) {
	var (
		ctx = context.Background()

		playerID = uuid.NewString()
		quest    = Quest{ID: uuid.NewString()}
	)

	t.Run("OK", func(t *testing.T) {
		getPlayerQuestProgression := BuildGetPlayerQuestProgression(func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
			return PlayerQuestProgression{Quest: quest, PlayerID: playerID}, nil
		})

		playerProgression, err := getPlayerQuestProgression(ctx, quest, playerID)
		assert.NoError(t, err)

		assert.Equal(t, playerID, playerProgression.PlayerID)
		assert.Equal(t, quest.ID, playerProgression.Quest.ID)
	})

	t.Run("Quest Not Found", func(t *testing.T) {
		getPlayerQuestProgression := BuildGetPlayerQuestProgression(func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
			return PlayerQuestProgression{}, ErrQuestNotFound
		})

		_, err := getPlayerQuestProgression(ctx, quest, playerID)
		assert.ErrorIs(t, err, ErrQuestNotFound)
	})

	t.Run("Random Error", func(t *testing.T) {
		getPlayerQuestProgression := BuildGetPlayerQuestProgression(func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
			return PlayerQuestProgression{}, errors.New("ant error")
		})

		_, err := getPlayerQuestProgression(ctx, quest, playerID)
		assert.Error(t, err)
	})
}

func TestBuildUpdatePlayerQuestProgressionFunc(t *testing.T) {
	var (
		ctx = context.Background()

		playerID = uuid.NewString()
	)

	t.Run("All Tasks Completed", func(t *testing.T) {
		quest := Quest{
			ID:     uuid.NewString(),
			GameID: uuid.NewString(),
			Tasks: []Task{
				{ID: uuid.NewString(), Rule: `{"==": [{"var": "fields.bool"}, true]}`},
			},
		}

		updatePlayerQuestProgressionFunc := BuildUpdatePlayerQuestProgressionFunc(
			func(ctx context.Context, progression PlayerQuestProgression) error {
				return nil
			},
			func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
				progression := make([]PlayerTaskProgression, len(quest.Tasks))
				for i, task := range quest.Tasks {
					progression[i] = PlayerTaskProgression{Task: task}
				}

				return PlayerQuestProgression{
					PlayerID:         playerID,
					Quest:            quest,
					TasksProgression: progression,
				}, nil
			},
			func(ctx context.Context, quest Quest, tasksCompleted []string, playerID string) (PlayerQuestProgression, error) {
				progression := make([]PlayerTaskProgression, len(tasksCompleted))
				for i, id := range tasksCompleted {
					progression[i] = PlayerTaskProgression{Task: Task{ID: id}, CompletedAt: time.Now()}
				}

				return PlayerQuestProgression{
					PlayerID:         playerID,
					Quest:            quest,
					TasksProgression: progression,
				}, nil
			},
		)

		progression, err := updatePlayerQuestProgressionFunc(ctx, quest, playerID, `{"fields": {"bool": true}}`)
		assert.NoError(t, err)

		assert.Equal(t, playerID, progression.PlayerID)
		assert.Equal(t, quest.ID, progression.Quest.ID)
		assert.Len(t, progression.TasksProgression, len(quest.Tasks))
		for _, task := range progression.TasksProgression {
			assert.NotEmpty(t, task.CompletedAt)
		}
	})

	t.Run("Nothing Completed", func(t *testing.T) {
		quest := Quest{
			ID:     uuid.NewString(),
			GameID: uuid.NewString(),
			Tasks: []Task{
				{ID: uuid.NewString(), Rule: `{"==": [{"var": "fields.bool"}, true]}`},
			},
		}

		updatePlayerQuestProgressionFunc := BuildUpdatePlayerQuestProgressionFunc(
			nil,
			func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
				progression := make([]PlayerTaskProgression, len(quest.Tasks))
				for i, task := range quest.Tasks {
					progression[i] = PlayerTaskProgression{Task: task}
				}

				return PlayerQuestProgression{
					PlayerID:         playerID,
					Quest:            quest,
					TasksProgression: progression,
				}, nil
			},
			nil,
		)

		progression, err := updatePlayerQuestProgressionFunc(ctx, quest, playerID, `{"fields": {"bool": false}}`)
		assert.NoError(t, err)

		assert.Equal(t, playerID, progression.PlayerID)
		assert.Equal(t, quest.ID, progression.Quest.ID)
		assert.Len(t, progression.TasksProgression, len(quest.Tasks))
		for _, task := range progression.TasksProgression {
			assert.Empty(t, task.CompletedAt)
		}
	})

	t.Run("Get Progression Error", func(t *testing.T) {
		quest := Quest{
			ID:     uuid.NewString(),
			GameID: uuid.NewString(),
			Tasks: []Task{
				{ID: uuid.NewString(), Rule: `{"==": [{"var": "fields.bool"}, true]}`},
			},
		}

		updatePlayerQuestProgressionFunc := BuildUpdatePlayerQuestProgressionFunc(
			nil,
			func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
				return PlayerQuestProgression{}, errors.New("any error")
			},
			nil,
		)

		progression, err := updatePlayerQuestProgressionFunc(ctx, quest, playerID, `{"fields": {"bool": false}}`)
		assert.Error(t, err)

		assert.Empty(t, progression.PlayerID)
		assert.Empty(t, progression.Quest.ID)
	})

	t.Run("Update Progression Error", func(t *testing.T) {
		quest := Quest{
			ID:     uuid.NewString(),
			GameID: uuid.NewString(),
			Tasks: []Task{
				{ID: uuid.NewString(), Rule: `{"==": [{"var": "fields.bool"}, true]}`},
			},
		}

		updatePlayerQuestProgressionFunc := BuildUpdatePlayerQuestProgressionFunc(
			nil,
			func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
				progression := make([]PlayerTaskProgression, len(quest.Tasks))
				for i, task := range quest.Tasks {
					progression[i] = PlayerTaskProgression{Task: task}
				}

				return PlayerQuestProgression{
					PlayerID:         playerID,
					Quest:            quest,
					TasksProgression: progression,
				}, nil
			},
			func(ctx context.Context, quest Quest, tasksCompleted []string, playerID string) (PlayerQuestProgression, error) {
				return PlayerQuestProgression{}, errors.New("any error")
			},
		)

		progression, err := updatePlayerQuestProgressionFunc(ctx, quest, playerID, `{"fields": {"bool": true}}`)
		assert.Error(t, err)

		assert.Empty(t, progression.PlayerID)
		assert.Empty(t, progression.Quest.ID)
	})

	t.Run("Progression Notifier Error", func(t *testing.T) {
		quest := Quest{
			ID:     uuid.NewString(),
			GameID: uuid.NewString(),
			Tasks: []Task{
				{ID: uuid.NewString(), Rule: `{"==": [{"var": "fields.bool"}, true]}`},
			},
		}

		updatePlayerQuestProgressionFunc := BuildUpdatePlayerQuestProgressionFunc(
			func(ctx context.Context, progression PlayerQuestProgression) error {
				return errors.New("any error")
			},
			func(ctx context.Context, quest Quest, playerID string) (PlayerQuestProgression, error) {
				progression := make([]PlayerTaskProgression, len(quest.Tasks))
				for i, task := range quest.Tasks {
					progression[i] = PlayerTaskProgression{Task: task}
				}

				return PlayerQuestProgression{
					PlayerID:         playerID,
					Quest:            quest,
					TasksProgression: progression,
				}, nil
			},
			func(ctx context.Context, quest Quest, tasksCompleted []string, playerID string) (PlayerQuestProgression, error) {
				progression := make([]PlayerTaskProgression, len(tasksCompleted))
				for i, id := range tasksCompleted {
					progression[i] = PlayerTaskProgression{Task: Task{ID: id}, CompletedAt: time.Now()}
				}

				return PlayerQuestProgression{
					PlayerID:         playerID,
					Quest:            quest,
					TasksProgression: progression,
				}, nil
			},
		)

		progression, err := updatePlayerQuestProgressionFunc(ctx, quest, playerID, `{"fields": {"bool": true}}`)
		assert.Error(t, err)

		assert.Empty(t, progression.PlayerID)
		assert.Empty(t, progression.Quest.ID)
	})
}
