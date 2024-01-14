package postgres

import (
	"context"

	"github.com/gabarcia/metagaming-api/internal/infra/storage/postgres/internal/sqlc"
	"github.com/gabarcia/metagaming-api/internal/quest"
	"github.com/google/uuid"
)

func sqlcQuestToDomain(q sqlc.Quest, ts []sqlc.Task) quest.Quest {
	tasks := make([]quest.Task, len(ts))
	for i, t := range ts {
		tasks[i] = sqlcTaskToDomain(t)
	}

	return quest.Quest{
		CreatedAt:   q.CreatedAt.Time,
		UpdatedAt:   q.UpdatedAt.Time,
		DeletedAt:   q.DeletedAt.Time,
		ID:          q.ID.String(),
		GameID:      q.GameID,
		Name:        q.Name,
		Description: q.Description,
		Tasks:       tasks,
	}
}

func (c connection) CreateQuest(ctx context.Context, data quest.NewQuestData) (quest.Quest, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return quest.Quest{}, err
	}
	defer tx.Rollback(context.Background())

	queries := c.queries.WithTx(tx)

	questData, err := queries.CreateQuest(ctx, sqlc.CreateQuestParams{
		GameID:      data.GameID,
		Name:        data.Name,
		Description: data.Description,
	})
	if err != nil {
		return quest.Quest{}, err
	}

	tasksData, err := createQuestTasks(ctx, queries, questData.ID, data.Tasks)
	if err != nil {
		return quest.Quest{}, err
	}

	return sqlcQuestToDomain(questData, tasksData), tx.Commit(ctx)
}

func (c connection) SoftDeleteQuestByIDAndGameID(ctx context.Context, id, gameID string) error {
	questID, err := uuid.Parse(id)
	if err != nil {
		return quest.ErrInvalidQuestID
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	queries := c.queries.WithTx(tx)

	affected, err := queries.SoftDeleteQuestByIDAndGameID(ctx, sqlc.SoftDeleteQuestByIDAndGameIDParams{
		ID:     questID,
		GameID: gameID,
	})
	if err != nil {
		return err
	}

	if affected == 0 {
		return quest.ErrQuestNotFound
	}

	if err = queries.SoftDeleteTasksByQuestID(ctx, questID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
