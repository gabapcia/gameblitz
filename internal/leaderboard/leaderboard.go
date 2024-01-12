package leaderboard

import (
	"context"
	"errors"
	"slices"
	"time"
)

var (
	ErrValidationError        = errors.New("validation error")
	ErrInvalidName            = errors.New("invalid name")
	ErrInvalidGameID          = errors.New("invalid game id")
	ErrInvalidStartDate       = errors.New("invalid start date")
	ErrInvalidAggregationMode = errors.New("invalid aggregation mode")
	ErrInvalidOrdering        = errors.New("invalid ordering")
	ErrEndDateBeforeStartDate = errors.New("end date must be after the start date")

	ErrInvalidLeaderboardID = errors.New("invalid leaderboard id")
	ErrLeaderboardNotFound  = errors.New("leaderboard not found")
)

const (
	AggregationModeInc = "INC"
	AggregationModeMax = "MAX"
	AggregationModeMin = "MIN"

	OrderingAsc  = "ASC"
	OrderingDesc = "DESC"
)

var (
	AggregationModes = []string{
		AggregationModeInc,
		AggregationModeMax,
		AggregationModeMin,
	}
	OrderingModes = []string{
		OrderingAsc,
		OrderingDesc,
	}
)

type NewLeaderboardData struct {
	GameID          string    // The ID from the game that is responsible for the leaderboard
	Name            string    // Leaderboard's name
	Description     string    // Leaderboard's description
	StartAt         time.Time // Time that the leaderboard should start working
	EndAt           time.Time // Time that the leaderboard will be closed for new updates
	AggregationMode string    // Data aggregation mode
	Ordering        string    // Leaderboard ranking order
}

type Leaderboard struct {
	CreatedAt       time.Time // Time that the leaderboard was created
	UpdatedAt       time.Time // Last time that the leaderboard info was updated
	DeletedAt       time.Time // Time that the leaderboard was deleted
	ID              string    // Leaderboard's ID
	GameID          string    // The ID from the game that is responsible for the leaderboard
	Name            string    // Leaderboard's name
	Description     string    // Leaderboard's description
	StartAt         time.Time // Time that the leaderboard should start working
	EndAt           time.Time // Time that the leaderboard will be closed for new updates
	AggregationMode string    // Data aggregation mode
	Ordering        string    // Leaderboard ranking order
}

func (l NewLeaderboardData) validate() error {
	errList := make([]error, 0)

	if l.Name == "" {
		errList = append(errList, ErrInvalidName)
	}

	if l.GameID == "" {
		errList = append(errList, ErrInvalidGameID)
	}

	if l.StartAt.IsZero() {
		errList = append(errList, ErrInvalidStartDate)
	}

	if !slices.Contains(AggregationModes, l.AggregationMode) {
		errList = append(errList, ErrInvalidAggregationMode)
	}

	if !slices.Contains(OrderingModes, l.Ordering) {
		errList = append(errList, ErrInvalidOrdering)
	}

	if !l.EndAt.IsZero() && l.EndAt.Before(l.StartAt) {
		errList = append(errList, ErrEndDateBeforeStartDate)
	}

	if len(errList) > 0 {
		errList = append(errList, ErrValidationError)
	}

	return errors.Join(errList...)
}

func (l Leaderboard) Closed() bool {
	now := time.Now()
	return !l.DeletedAt.IsZero() || now.Before(l.StartAt) || (!l.EndAt.IsZero() && now.After(l.EndAt))
}

func BuildCreateFunc(storageCreateFunc StorageCreateLeaderboardFunc) CreateFunc {
	return func(ctx context.Context, data NewLeaderboardData) (Leaderboard, error) {
		if err := data.validate(); err != nil {
			return Leaderboard{}, err
		}

		return storageCreateFunc(ctx, data)
	}
}

func BuildGetByIDAndGameIDFunc(storageGetFunc StorageGetLeaderboardByIDAndGameIDFunc) GetByIDAndGameIDFunc {
	return func(ctx context.Context, id, gameID string) (Leaderboard, error) {
		return storageGetFunc(ctx, id, gameID)
	}
}

func BuildSoftDeleteFunc(storageSoftDeleteFunc StorageSoftDeleteLeaderboardFunc) SoftDeleteFunc {
	return func(ctx context.Context, id, gameID string) error {
		return storageSoftDeleteFunc(ctx, id, gameID)
	}
}
