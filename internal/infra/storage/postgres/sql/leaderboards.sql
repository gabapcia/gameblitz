-- name: CreateLeaderboard :one
INSERT INTO "leaderboards" ("game_id", "name", "description", "start_at", "end_at", "aggregation_mode", "data_type", "ordering")
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING "id";
