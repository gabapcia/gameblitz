-- name: CreateLeaderboard :one
INSERT INTO "leaderboards" ("game_id", "name", "description", "start_at", "end_at", "aggregation_mode", "data_type", "ordering")
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetLeaderboardByIDAndGameID :one
SELECT *
FROM "leaderboards" l
WHERE l."id" = $1 AND l."game_id" = $2 AND l."deleted_at" IS NULL;

-- name: SoftDeleteLeaderboard :execrows
UPDATE "leaderboards"
SET
    "deleted_at" = NOW()
WHERE "id" = $1 AND "game_id" = $2 AND "deleted_at" IS NULL;
