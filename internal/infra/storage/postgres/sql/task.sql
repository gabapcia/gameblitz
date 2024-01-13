-- name: CreateTask :one
INSERT INTO "tasks" ("quest_id", "name", "description", "depends_on", "rule")
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
