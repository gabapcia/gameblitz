-- name: CreateQuest :one
INSERT INTO "quests" ("game_id", "name", "description")
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetQuestByIDAndGameID :one
SELECT *
FROM "quests" q
WHERE q."id" = $1 AND q."game_id" = $2 AND q."deleted_at" IS NULL
LIMIT 1;

-- name: SoftDeleteQuestByIDAndGameID :execrows
UPDATE "quests"
SET
    "deleted_at" = NOW()
WHERE "id" = $1 AND "game_id" = $2 AND "deleted_at" IS NULL;
