DROP INDEX IF EXISTS "idx_leaderboard_created_at" CASCADE;
DROP INDEX IF EXISTS "idx_leaderboard_deleted_at" CASCADE;
DROP INDEX IF EXISTS "idx_leaderboard_game_id" CASCADE;

DROP TABLE IF EXISTS "leaderboards" CASCADE;
