CREATE TABLE IF NOT EXISTS "quests" (
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "game_id" VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',

    CONSTRAINT "name_not_empty_check" CHECK (TRIM("name") <> ''),
    CONSTRAINT "game_id_not_empty_check" CHECK (TRIM("game_id") <> '')
);

CREATE INDEX IF NOT EXISTS "idx_quest_created_at" ON "quests" ("created_at");
CREATE INDEX IF NOT EXISTS "idx_quest_deleted_at" ON "quests" ("deleted_at");
CREATE INDEX IF NOT EXISTS "idx_quest_game_id" ON "quests" ("game_id");

CREATE TABLE IF NOT EXISTS "tasks" (
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "quest_id" UUID NOT NULL REFERENCES "quests"("id") ON DELETE CASCADE,
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "required_for_completion" BOOLEAN NOT NULL,
    "rule" TEXT NOT NULL,

    CONSTRAINT "name_not_empty_check" CHECK (TRIM("name") <> '')
);

CREATE INDEX IF NOT EXISTS "idx_task_created_at" ON "tasks" ("created_at");
CREATE INDEX IF NOT EXISTS "idx_task_deleted_at" ON "tasks" ("deleted_at");
CREATE INDEX IF NOT EXISTS "idx_task_quest_id" ON "tasks" ("quest_id");
CREATE INDEX IF NOT EXISTS "idx_task_required_for_completion" ON "tasks" ("required_for_completion");

CREATE TABLE IF NOT EXISTS "tasks_dependencies" (
    "this_task" UUID NOT NULL REFERENCES "tasks" ("id") ON DELETE CASCADE,
    "depends_on_task" UUID NOT NULL REFERENCES "tasks" ("id") ON DELETE CASCADE,

    PRIMARY KEY ("this_task", "depends_on_task")
);

CREATE INDEX IF NOT EXISTS "idx_tasks_dependency_this_task" ON "tasks_dependencies" ("this_task");
CREATE INDEX IF NOT EXISTS "idx_tasks_dependency_depends_on_task" ON "tasks_dependencies" ("depends_on_task");

CREATE OR REPLACE VIEW "tasks_with_its_dependencies" AS
    SELECT t.*, ARRAY_REMOVE(ARRAY_AGG(td."depends_on_task"), NULL)::UUID[] AS "depends_on" 
    FROM "tasks" t
    LEFT JOIN "tasks_dependencies" td on t."id" = td."this_task"
    GROUP BY t."id"
    ORDER BY t."created_at" ASC;
