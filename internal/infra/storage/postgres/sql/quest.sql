-- name: CreateQuest :one
INSERT INTO "quests" ("game_id", "name", "description")
VALUES ($1, $2, $3)
RETURNING *;

-- name: SoftDeleteQuestByIDAndGameID :execrows
UPDATE "quests" q
SET
    q."deleted_at" = NOW()
WHERE q."id" = $1 AND q."game_id" = $2 AND q."deleted_at" IS NULL;
