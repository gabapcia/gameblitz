-- name: CreateTask :one
INSERT INTO "tasks" ("quest_id", "name", "description", "required_for_completion", "rule")
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: RegisterTaskDependency :exec
INSERT INTO "tasks_dependencies" ("this_task", "depends_on_task")
VALUES ($1, $2);

-- name: ListTasksByQuestID :many
SELECT *
FROM "tasks_with_its_dependencies" t
WHERE t."quest_id" = $1 AND t."deleted_at" IS NULL;

-- name: SoftDeleteTasksByQuestID :exec
UPDATE "tasks"
SET
    "deleted_at" = NOW()
WHERE "quest_id" = $1 AND "deleted_at" IS NULL;
