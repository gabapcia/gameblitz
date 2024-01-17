-- name: StartPlayerTasksForQuest :many
WITH "player_quest_tasks_created" AS (
    INSERT INTO "player_quest_tasks" ("player_id", "task_id")
    SELECT $1, t."id"
    FROM "tasks" t
    WHERE t."quest_id" = $2 AND t."depends_on" IS NULL
    RETURNING *
)
SELECT pqt.*, sqlc.embed(t), ARRAY_AGG(ptd."depends_on_task")::UUID[] AS "depends_on"
FROM "player_quest_tasks_created" pqt
JOIN "tasks" t ON t."id" = pqt."task_id"
LEFT JOIN "tasks_dependencies" ptd ON t."id" = ptd."this_task"
GROUP BY pqt."task_id";

-- name: GetPlayerQuestTasks :many
SELECT pqt.*, sqlc.embed(t), ARRAY_AGG(ptd."depends_on_task")::UUID[] AS "depends_on"
FROM "player_quest_tasks" pqt
JOIN "tasks" t ON t."id" = pqt."task_id"
LEFT JOIN "tasks_dependencies" ptd ON t."id" = ptd."this_task"
WHERE pqt."player_id" = $1 AND t."quest_id" = $2
GROUP BY pqt."task_id";
