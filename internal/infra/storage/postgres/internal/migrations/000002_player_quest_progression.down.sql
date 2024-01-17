DROP INDEX IF EXISTS "idx_player_quest_task_player_quest_id" CASCADE;

DROP TABLE IF EXISTS "player_quest_tasks" CASCADE;

DROP FUNCTION IF EXISTS "validate_task_belongs_to_quest" CASCADE;

DROP TABLE IF EXISTS "player_quests" CASCADE;
