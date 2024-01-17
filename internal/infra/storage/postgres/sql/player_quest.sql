------------------------
-- Start Player Quest --
------------------------

-- name: StartPlayerQuest :one
INSERT INTO "player_quests" ("player_id", "quest_id")
SELECT $1, q."id"
FROM "quests" q
WHERE q."id" = sqlc.arg('quest_id') AND q."deleted_at" IS NULL
RETURNING *;

-- name: StartPlayerTasksForQuest :many
WITH "player_quest_tasks_created" AS (
    INSERT INTO "player_quest_tasks" ("player_id", "player_quest_id", "task_id")
    SELECT pq."player_id", $1, t."id"
    FROM "player_quests" pq
    JOIN "tasks_with_its_dependencies" t ON t."quest_id" = pq."quest_id"
    WHERE pq."id" = $1 AND ARRAY_LENGTH(t."depends_on", 1) IS NULL
    RETURNING *
)
SELECT pqt.*, sqlc.embed(twd)
FROM "player_quest_tasks_created" pqt
JOIN "tasks_with_its_dependencies" twd ON twd."id" = pqt."task_id";

-----------------------
-- Get Player Quests --
-----------------------

-- name: GetPlayerQuest :one
SELECT *
FROM "player_quests" pq
WHERE pq."player_id" = $1 AND pq."quest_id" = $2;

-- name: GetPlayerQuestTasks :many
SELECT pqt.*, sqlc.embed(t)
FROM "player_quest_tasks" pqt
JOIN "tasks_with_its_dependencies" t ON t."id" = pqt."task_id"
WHERE pqt."player_id" = $1 AND t."quest_id" = $2;
