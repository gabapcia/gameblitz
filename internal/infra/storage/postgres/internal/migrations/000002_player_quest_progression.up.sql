CREATE TABLE IF NOT EXISTS "player_quests" (
    "started_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "player_id" VARCHAR NOT NULL,
    "quest_id" UUID NOT NULL REFERENCES "quests" ("id") ON DELETE CASCADE,
    "completed_at" TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT "player_quest_unique" UNIQUE ("player_id", "quest_id")
);

CREATE TABLE IF NOT EXISTS "player_quest_tasks" (
    "started_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "player_id" VARCHAR NOT NULL,
    "task_id" UUID NOT NULL REFERENCES "tasks" ("id") ON DELETE CASCADE,
    "completed_at" TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT "player_quest_task_unique" UNIQUE ("player_id", "task_id")
);
