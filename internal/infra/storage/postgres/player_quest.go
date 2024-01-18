package postgres

import (
	"context"
	"errors"

	"github.com/gabarcia/metagaming-api/internal/infra/storage/postgres/internal/sqlc"
	"github.com/gabarcia/metagaming-api/internal/quest"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func sqcStartQuestForPlayerDataToDomain(pq sqlc.PlayerQuest, q quest.Quest, ts []sqlc.StartPlayerTasksForQuestRow) quest.PlayerQuestProgression {
	tasksProgression := make([]quest.PlayerTaskProgression, len(ts))
	for i, t := range ts {
		tasksProgression[i] = quest.PlayerTaskProgression{
			StartedAt:   t.StartedAt.Time,
			UpdatedAt:   t.UpdatedAt.Time,
			Task:        sqlcTaskWithItsDependenciesToDomain(t.TasksWithItsDependency),
			CompletedAt: t.CompletedAt.Time,
		}
	}

	return quest.PlayerQuestProgression{
		StartedAt:        pq.StartedAt.Time,
		UpdatedAt:        pq.UpdatedAt.Time,
		PlayerID:         pq.PlayerID,
		Quest:            q,
		CompletedAt:      pq.CompletedAt.Time,
		TasksProgression: tasksProgression,
	}
}

func sqlcGetPlayerQuestDataToDomain(pq sqlc.PlayerQuest, q quest.Quest, ts []sqlc.GetPlayerQuestTasksRow) quest.PlayerQuestProgression {
	tasksProgression := make([]quest.PlayerTaskProgression, len(ts))
	for i, t := range ts {
		tasksProgression[i] = quest.PlayerTaskProgression{
			StartedAt:   t.StartedAt.Time,
			UpdatedAt:   t.UpdatedAt.Time,
			Task:        sqlcTaskWithItsDependenciesToDomain(t.TasksWithItsDependency),
			CompletedAt: t.CompletedAt.Time,
		}
	}

	return quest.PlayerQuestProgression{
		StartedAt:        pq.StartedAt.Time,
		UpdatedAt:        pq.UpdatedAt.Time,
		PlayerID:         pq.PlayerID,
		Quest:            q,
		CompletedAt:      pq.CompletedAt.Time,
		TasksProgression: tasksProgression,
	}
}

func (c connection) StartQuestForPlayer(ctx context.Context, q quest.Quest, playerID string) (quest.PlayerQuestProgression, error) {
	questID, err := uuid.Parse(q.ID)
	if err != nil {
		return quest.PlayerQuestProgression{}, quest.ErrInvalidQuestID
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return quest.PlayerQuestProgression{}, err
	}
	defer tx.Rollback(context.Background())

	queries := c.queries.WithTx(tx)

	playerQuestData, err := queries.StartPlayerQuest(ctx, sqlc.StartPlayerQuestParams{
		PlayerID: playerID,
		QuestID:  questID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = quest.ErrQuestNotFound
		}

		return quest.PlayerQuestProgression{}, err
	}

	playerQuestTasksData, err := queries.StartPlayerTasksForQuest(ctx, playerQuestData.ID)
	if err != nil {
		return quest.PlayerQuestProgression{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return quest.PlayerQuestProgression{}, err
	}

	return sqcStartQuestForPlayerDataToDomain(playerQuestData, q, playerQuestTasksData), nil
}

func (c connection) GetPlayerQuestProgression(ctx context.Context, q quest.Quest, playerID string) (quest.PlayerQuestProgression, error) {
	questID, err := uuid.Parse(q.ID)
	if err != nil {
		return quest.PlayerQuestProgression{}, quest.ErrInvalidQuestID
	}

	playerQuestData, err := c.queries.GetPlayerQuest(ctx, sqlc.GetPlayerQuestParams{
		PlayerID: playerID,
		QuestID:  questID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = quest.ErrQuestNotFound
		}

		return quest.PlayerQuestProgression{}, err
	}

	playerTasksData, err := c.queries.GetPlayerQuestTasks(ctx, sqlc.GetPlayerQuestTasksParams{
		PlayerID: playerID,
		QuestID:  playerQuestData.QuestID,
	})
	if err != nil {
		return quest.PlayerQuestProgression{}, err
	}

	return sqlcGetPlayerQuestDataToDomain(playerQuestData, q, playerTasksData), nil
}

func (c connection) UpdatePlayerQuestProgression(ctx context.Context, q quest.Quest, tc []string, playerID string) (quest.PlayerQuestProgression, error) {
	questID, err := uuid.Parse(q.ID)
	if err != nil {
		return quest.PlayerQuestProgression{}, quest.ErrInvalidQuestID
	}

	tasksCompleted := make([]uuid.UUID, len(tc))
	for i, taskIDRaw := range tc {
		taskID, err := uuid.Parse(taskIDRaw)
		if err != nil {
			return quest.PlayerQuestProgression{}, quest.ErrInvalidTaskID
		}

		tasksCompleted[i] = taskID
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return quest.PlayerQuestProgression{}, err
	}
	defer tx.Rollback(context.Background())

	queries := c.queries.WithTx(tx)

	err = queries.MarkPlayerQuestTasksAsCompleted(ctx, sqlc.MarkPlayerQuestTasksAsCompletedParams{
		PlayerID:       playerID,
		TasksCompleted: tasksCompleted,
	})
	if err != nil {
		return quest.PlayerQuestProgression{}, err
	}

	err = queries.StartPlayerTasksThatHadTheDependenciesCompleted(ctx, sqlc.StartPlayerTasksThatHadTheDependenciesCompletedParams{
		QuestID:  questID,
		PlayerID: playerID,
	})
	if err != nil {
		return quest.PlayerQuestProgression{}, err
	}

	err = queries.MarkPlayerQuestAsCompleted(ctx, sqlc.MarkPlayerQuestAsCompletedParams{
		QuestID:  questID,
		PlayerID: playerID,
	})
	if err != nil {
		return quest.PlayerQuestProgression{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return quest.PlayerQuestProgression{}, err
	}

	return c.GetPlayerQuestProgression(ctx, q, playerID)
}
