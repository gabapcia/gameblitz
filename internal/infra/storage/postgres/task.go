package postgres

import (
	"context"

	"github.com/gabarcia/metagaming-api/internal/infra/storage/postgres/internal/sqlc"
	"github.com/gabarcia/metagaming-api/internal/quest"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/google/uuid"
)

func sqlcTaskToDomain(t sqlc.Task) quest.Task {
	var (
		dependsOnValue, _ = t.DependsOn.Value()
		dependsOn         = ""
	)
	if dependsOnValue != nil {
		dependsOn = dependsOnValue.(string)
	}

	return quest.Task{
		CreatedAt:   t.CreatedAt.Time,
		UpdatedAt:   t.UpdatedAt.Time,
		DeletedAt:   t.DeletedAt.Time,
		ID:          t.ID.String(),
		Name:        t.Name,
		Description: t.Description,
		DependsOn:   dependsOn,
		Rule:        t.Rule,
	}
}

func createQuestTasks(ctx context.Context, queries *sqlc.Queries, questID uuid.UUID, tasks []quest.NewTaskData) ([]sqlc.Task, error) {
	var (
		tasksCreated = make([]sqlc.Task, len(tasks))
		completed    = false
	)
	for !completed {
		completed = true
		for i, task := range tasks {
			if tasksCreated[i].ID.String() != uuid.Nil.String() {
				continue
			}

			var dependsOn pgtype.UUID
			if task.DependsOn != nil {
				dependsOnTaskID := tasksCreated[*task.DependsOn].ID.String()
				if dependsOnTaskID == uuid.Nil.String() {
					completed = false
					continue
				}

				if err := dependsOn.Scan(dependsOnTaskID); err != nil {
					return nil, err
				}
			}

			taskData, err := queries.CreateTask(ctx, sqlc.CreateTaskParams{
				QuestID:     questID,
				Name:        task.Name,
				Description: task.Description,
				DependsOn:   dependsOn,
				Rule:        task.Rule,
			})
			if err != nil {
				return nil, err
			}

			tasksCreated[i] = taskData
		}
	}

	return tasksCreated, nil
}
