package postgres

import (
	"context"

	"github.com/gabarcia/metagaming-api/internal/infra/storage/postgres/internal/sqlc"
	"github.com/gabarcia/metagaming-api/internal/quest"

	"github.com/google/uuid"
)

func sqlcTaskToDomain(t sqlc.Task, dependsOnUUIDs []uuid.UUID) quest.Task {
	dependsOn := make([]string, len(dependsOnUUIDs))
	for _, td := range dependsOnUUIDs {
		dependsOn = append(dependsOn, td.String())
	}

	return quest.Task{
		CreatedAt:             t.CreatedAt.Time,
		UpdatedAt:             t.UpdatedAt.Time,
		DeletedAt:             t.DeletedAt.Time,
		ID:                    t.ID.String(),
		Name:                  t.Name,
		Description:           t.Description,
		DependsOn:             dependsOn,
		RequiredForCompletion: t.RequiredForCompletion,
		Rule:                  t.Rule,
	}
}

func createQuestTasks(ctx context.Context, queries *sqlc.Queries, questID uuid.UUID, tasks []quest.NewTaskData) (map[sqlc.Task][]uuid.UUID, error) {
	var (
		rawDependenciesMap = make(map[uuid.UUID][]int)
		tasksCreatedRows   = make([]sqlc.Task, len(tasks))
	)
	for i, task := range tasks {
		taskData, err := queries.CreateTask(ctx, sqlc.CreateTaskParams{
			QuestID:               questID,
			Name:                  task.Name,
			Description:           task.Description,
			RequiredForCompletion: task.RequiredForCompletion,
			Rule:                  task.Rule,
		})
		if err != nil {
			return nil, err
		}

		tasksCreatedRows[i] = taskData
		rawDependenciesMap[taskData.ID] = task.DependsOn
	}

	tasksCreated := make(map[sqlc.Task][]uuid.UUID)
	for _, taskData := range tasksCreatedRows {
		dependsOn := make([]uuid.UUID, len(rawDependenciesMap[taskData.ID]))
		for i, dependsOnIndex := range rawDependenciesMap[taskData.ID] {
			err := queries.RegisterTaskDependency(ctx, sqlc.RegisterTaskDependencyParams{
				ThisTask:      taskData.ID,
				DependsOnTask: tasksCreatedRows[dependsOnIndex].ID,
			})
			if err != nil {
				return nil, err
			}

			dependsOn[i] = tasksCreatedRows[dependsOnIndex].ID
		}

		tasksCreated[taskData] = dependsOn
	}

	return tasksCreated, nil
}
