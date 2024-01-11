CREATE TABLE IF NOT EXISTS "leaderboards" (
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMPTZ DEFAULT NULL,
    "id" UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "game_id" VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "start_at" TIMESTAMPTZ NOT NULL,
    "end_at" TIMESTAMPTZ DEFAULT NULL,
    "aggregation_mode" VARCHAR NOT NULL,
    "data_type" VARCHAR NOT NULL,
    "ordering" VARCHAR NOT NULL,

    CONSTRAINT "game_id_check" CHECK (TRIM("game_id") <> ''),
    CONSTRAINT "name_check" CHECK (TRIM("name") <> ''),
    CONSTRAINT "end_date_after_start_date_check" CHECK ("end_at" IS NULL OR "end_at" > "start_at")
);

CREATE INDEX IF NOT EXISTS "idx_leaderboard_created_at" ON "leaderboards" ("created_at");
CREATE INDEX IF NOT EXISTS "idx_leaderboard_deleted_at" ON "leaderboards" ("deleted_at");
CREATE INDEX IF NOT EXISTS "idx_leaderboard_game_id" ON "leaderboards" ("game_id");
