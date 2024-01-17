DROP VIEW IF EXISTS "tasks_with_its_dependencies";

DROP INDEX IF EXISTS "idx_tasks_dependency_this_task" CASCADE;
DROP INDEX IF EXISTS "idx_tasks_dependency_depends_on_task" CASCADE;

DROP TABLE IF EXISTS "tasks_dependencies" CASCADE;

DROP INDEX IF EXISTS "idx_task_created_at" CASCADE;
DROP INDEX IF EXISTS "idx_task_deleted_at" CASCADE;
DROP INDEX IF EXISTS "idx_task_quest_id" CASCADE;
DROP INDEX IF EXISTS "idx_task_depends_on" CASCADE;
DROP INDEX IF EXISTS "idx_task_required_for_completion" CASCADE;

DROP TABLE IF EXISTS "tasks" CASCADE;

DROP INDEX IF EXISTS "idx_quest_created_at" CASCADE;
DROP INDEX IF EXISTS "idx_quest_deleted_at" CASCADE;
DROP INDEX IF EXISTS "idx_quest_game_id" CASCADE;

DROP TABLE IF EXISTS "quests" CASCADE;
