-- name: StartPlayerQuest :one
INSERT INTO "player_quests" ("player_id", "quest_id")
SELECT $1, q."id"
FROM "quests" q
WHERE q."id" = sqlc.arg('quest_id') AND q."deleted_at" IS NULL
RETURNING *;

-- name: GetPlayerQuest :one
SELECT *
FROM "player_quests" pq
WHERE pq."player_id" = $1 AND pq."quest_id" = $2;
