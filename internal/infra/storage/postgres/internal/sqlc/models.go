// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package sqlc

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Leaderboard struct {
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
	DeletedAt       pgtype.Timestamptz
	ID              uuid.UUID
	GameID          string
	Name            string
	Description     string
	StartAt         pgtype.Timestamptz
	EndAt           pgtype.Timestamptz
	AggregationMode string
	Ordering        string
}
