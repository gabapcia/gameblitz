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

---------------------------------------
-- Mark Quest And Tasks As Completed --
---------------------------------------

-- name: MarkPlayerQuestTasksAsCompleted :exec
UPDATE "player_quest_tasks"
SET
    "updated_at" = NOW(),
    "completed_at" = NOW()
WHERE
    "player_id" = $1 AND
    "completed_at" IS NULL AND
    "task_id" = ANY(sqlc.arg('tasks_completed')::UUID[]);

-- name: StartPlayerTasksThatHadTheDependenciesCompleted :exec
WITH "pq_tasks_status" AS (
    SELECT pqt."task_id", (pqt."completed_at" IS NOT NULL) AS "completed"
    FROM "player_quests" pq
    JOIN "player_quest_tasks" pqt
        ON pqt."player_quest_id" = pq."id"
    WHERE
        pq."quest_id" = $1 AND pq."player_id" = $2
), "pq_pending_tasks" AS (
	SELECT t."id"
	FROM "tasks" t
	WHERE
	    t."quest_id" = $1 AND
	    t."id" NOT IN (SELECT "task_id" FROM "pq_tasks_status")
), "pq_tasks_ready_to_start" AS (
    SELECT td."this_task" AS "id"
    FROM "tasks_dependencies" td
    WHERE
        td."this_task" IN (SELECT "id" FROM "pq_pending_tasks")
    GROUP BY td."this_task"
    HAVING
        ARRAY_AGG(td."depends_on_task") <@ (
            SELECT ARRAY_AGG("task_id")
            FROM "pq_tasks_status"
            WHERE "completed" = TRUE
        )
)
INSERT INTO "player_quest_tasks" ("player_id", "player_quest_id", "task_id")
SELECT $2, pq2."id", trs."id"
FROM "pq_tasks_ready_to_start" trs
CROSS JOIN "player_quests" pq2
WHERE
    pq2."quest_id" = $1 AND pq2."player_id" = $2;

-- name: MarkPlayerQuestAsCompleted :exec
WITH "completion_list" AS (
	SELECT (pqt."completed_at" IS NOT NULL) AS "completed"
	FROM "tasks" t
	LEFT JOIN "player_quest_tasks" pqt
        ON t.id = pqt."task_id" AND pqt."player_id" = $2
	WHERE
        t."quest_id" = $1 AND
        t."required_for_completion" = TRUE
)
UPDATE "player_quests"
SET
    "updated_at" = NOW(),
    "completed_at" = NOW()
WHERE 
	"player_quests"."player_id" = $2 AND
	"player_quests"."quest_id" = $1 AND 
	"player_quests"."completed_at" IS NULL AND 
	TRUE = ALL((SELECT "completed" FROM "completion_list"));
