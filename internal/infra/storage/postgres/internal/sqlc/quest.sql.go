// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: quest.sql

package sqlc

import (
	"context"
)

const createQuest = `-- name: CreateQuest :one
INSERT INTO "quests" ("game_id", "name", "description")
VALUES ($1, $2, $3)
RETURNING created_at, updated_at, deleted_at, id, game_id, name, description
`

type CreateQuestParams struct {
	GameID      string
	Name        string
	Description string
}

func (q *Queries) CreateQuest(ctx context.Context, arg CreateQuestParams) (Quest, error) {
	row := q.db.QueryRow(ctx, createQuest, arg.GameID, arg.Name, arg.Description)
	var i Quest
	err := row.Scan(
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.ID,
		&i.GameID,
		&i.Name,
		&i.Description,
	)
	return i, err
}
