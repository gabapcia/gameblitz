DROP INDEX IF EXISTS "idx_task_created_at" CASCADE;
DROP INDEX IF EXISTS "idx_task_deleted_at" CASCADE;
DROP INDEX IF EXISTS "idx_task_quest_id" CASCADE;
DROP INDEX IF EXISTS "idx_task_depends_on" CASCADE;

DROP TABLE IF EXISTS "tasks" CASCADE;

DROP INDEX IF EXISTS "idx_quest_created_at" CASCADE;
DROP INDEX IF EXISTS "idx_quest_deleted_at" CASCADE;
DROP INDEX IF EXISTS "idx_quest_game_id" CASCADE;

DROP TABLE IF EXISTS "quests" CASCADE;
