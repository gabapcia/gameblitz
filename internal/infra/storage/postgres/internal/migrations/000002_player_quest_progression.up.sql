CREATE TABLE IF NOT EXISTS "player_quests" (
    "started_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "player_id" VARCHAR NOT NULL,
    "quest_id" UUID NOT NULL REFERENCES "quests" ("id") ON DELETE CASCADE,
    "completed_at" TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT "player_quest_unique" UNIQUE ("player_id", "quest_id")
);

CREATE OR REPLACE FUNCTION validate_task_belongs_to_quest("task_id" UUID, "player_quest_id" UUID) RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1
        FROM "tasks" t
        JOIN "quests" q ON q."id" = t."quest_id"
        WHERE
            t."id" = "task_id" AND
            q."id" = (SELECT "quest_id" FROM "player_quests" WHERE "id" = "player_quest_id")
    );
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION validate_task_dependencies_completed("task_id_to_start" UUID, "pq_id" UUID) RETURNS BOOLEAN AS $$
BEGIN
    RETURN TRUE = ALL(
        SELECT (pqt."completed_at" IS NOT NULL)
        FROM "tasks_dependencies" td
        LEFT JOIN "player_quest_tasks" pqt ON 
            pqt."task_id" = td."depends_on_task" AND
            pqt."player_quest_id" = "pq_id"
        WHERE td."this_task" = "task_id_to_start"
    );
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS "player_quest_tasks" (
    "started_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "player_id" VARCHAR NOT NULL,
    "player_quest_id" UUID NOT NULL REFERENCES "player_quests" ("id") ON DELETE CASCADE,
    "task_id" UUID NOT NULL REFERENCES "tasks" ("id") ON DELETE CASCADE,
    "completed_at" TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT "player_quest_task_unique" UNIQUE ("player_id", "task_id"),
    CONSTRAINT "task_belongs_to_quest_check" CHECK (validate_task_belongs_to_quest("task_id", "player_quest_id")),
    CONSTRAINT "task_can_be_started_check" CHECK (validate_task_dependencies_completed("task_id", "player_quest_id"))
);

CREATE INDEX IF NOT EXISTS "idx_player_quest_task_player_quest_id" ON "player_quest_tasks" ("player_quest_id");
