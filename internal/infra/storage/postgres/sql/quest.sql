-- name: CreateQuest :one
INSERT INTO "quests" ("game_id", "name", "description")
VALUES ($1, $2, $3)
RETURNING *;
