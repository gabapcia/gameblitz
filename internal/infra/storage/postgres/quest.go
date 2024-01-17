package postgres

import (
	"context"
	"errors"

	"github.com/gabarcia/metagaming-api/internal/infra/storage/postgres/internal/sqlc"
	"github.com/gabarcia/metagaming-api/internal/quest"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func sqlcQuestToDomain(q sqlc.Quest, ts map[sqlc.Task][]uuid.UUID) quest.Quest {
	tasks := make([]quest.Task, 0)
	for t, ds := range ts {
		tasks = append(tasks, sqlcTaskToDomain(t, ds))
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

func (c connection) GetQuestByIDAndGameID(ctx context.Context, id, gameID string) (quest.Quest, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return quest.Quest{}, quest.ErrInvalidQuestID
	}

	questData, err := c.queries.GetQuestByIDAndGameID(ctx, sqlc.GetQuestByIDAndGameIDParams{
		ID:     uid,
		GameID: gameID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = quest.ErrQuestNotFound
		}

		return quest.Quest{}, err
	}

	tasksDataRaw, err := c.queries.ListTasksByQuestID(ctx, questData.ID)
	if err != nil {
		return quest.Quest{}, err
	}

	tasksData := make(map[sqlc.Task][]uuid.UUID)
	for _, row := range tasksDataRaw {
		tasksData[row.Task] = row.DependsOn
	}

	return sqlcQuestToDomain(questData, tasksData), nil
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
