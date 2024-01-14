-- name: CreateTask :one
INSERT INTO "tasks" ("quest_id", "name", "description", "depends_on", "rule")
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListTasksByQuestID :many
SELECT *
FROM "tasks" t
WHERE t."quest_id" = $1 AND t."deleted_at" IS NULL;

-- name: SoftDeleteTasksByQuestID :exec
UPDATE "tasks"
SET
    "deleted_at" = NOW()
WHERE "quest_id" = $1 AND "deleted_at" IS NULL;
