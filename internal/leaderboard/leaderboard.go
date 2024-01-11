package leaderboard

import (
	"errors"
	"slices"
	"time"
)

var (
	ErrValidationError        = errors.New("Validation error")
	ErrInvalidLeaderboard     = errors.New("Invalid leaderboard")
	ErrInvalidName            = errors.New("Invalid name")
	ErrInvalidGameID          = errors.New("Invalid game id")
	ErrInvalidStartDate       = errors.New("Invalid start date")
	ErrInvalidAggregationMode = errors.New("Invalid aggregation mode")
	ErrInvalidDataType        = errors.New("Invalid data type")
	ErrInvalidOrdering        = errors.New("Invalid ordering")
	ErrEndDateBeforeStartDate = errors.New("End date must be after the start date")
)

const (
	AggregationModeInc = "INC"
	AggregationModeMax = "MAX"
	AggregationModeMin = "MIN"

	DataTypeInt = "INT"

	OrderingAsc  = "ASC"
	OrderingDesc = "DESC"
)

var (
	AggregationModes = []string{
		AggregationModeInc,
		AggregationModeMax,
		AggregationModeMin,
	}
	DataTypes = []string{
		DataTypeInt,
	}
	OrderingModes = []string{
		OrderingAsc,
		OrderingDesc,
	}
)

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
	DataType        string    // Data type that the leaderboard should accept
	Ordering        string    // Leaderboard ranking order
}

func (l Leaderboard) validate() error {
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

	if !slices.Contains(DataTypes, l.DataType) {
		errList = append(errList, ErrInvalidDataType)
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
