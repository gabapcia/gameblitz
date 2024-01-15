package statistic

import (
	"context"
	"errors"
	"slices"
	"time"
)

var (
	ErrStatisticValidation    = errors.New("invalid statistic")
	ErrInvalidStatisticID     = errors.New("invalid id")
	ErrInvalidName            = errors.New("invalid name")
	ErrMissingGameID          = errors.New("missing game id")
	ErrInvalidAggregationMode = errors.New("invalid aggregation mode")
	ErrStatisticNotFound      = errors.New("statistic not found")
)

const (
	AggregationModeSum = "SUM"
	AggregationModeSub = "SUB"
	AggregationModeMax = "MAX"
	AggregationModeMin = "MIN"
)

var AggregationModes = []string{
	AggregationModeSum,
	AggregationModeSub,
	AggregationModeMax,
	AggregationModeMin,
}

type NewStatisticData struct {
	GameID          string    // ID of the game responsible for the statistic
	Name            string    // Statistic name
	Description     string    // Statistic details
	AggregationMode string    // Data aggregation mode
	InitialValue    *float64  // Initial statistic value for players
	Goal            *float64  // Goal value. nil means no goal
	Landmarks       []float64 // Statistic landmarks
}

type Statistic struct {
	CreatedAt       time.Time // Time that the statistic was created
	UpdatedAt       time.Time // Last time that the statistic was updated
	DeletedAt       time.Time // Time that the statistic was deleted
	ID              string    // Statistic ID
	GameID          string    // ID of the game responsible for the statistic
	Name            string    // Statistic name
	Description     string    // Statistic details
	AggregationMode string    // Data aggregation mode
	InitialValue    *float64  // Initial statistic value for players
	Goal            *float64  // Goal value. nil means no goal
	Landmarks       []float64 // Statistic landmarks
}

func (s NewStatisticData) validate() error {
	errList := make([]error, 0)

	if s.GameID == "" {
		errList = append(errList, ErrMissingGameID)
	}

	if s.Name == "" {
		errList = append(errList, ErrInvalidName)
	}

	if !slices.Contains(AggregationModes, s.AggregationMode) {
		errList = append(errList, ErrInvalidAggregationMode)
	}

	if len(errList) > 0 {
		errList = append(errList, ErrStatisticValidation)
	}

	return errors.Join(errList...)
}

func BuildCreateStatisticFunc(storageCreateStatisticFunc StorageCreateStatisticFunc) CreateFunc {
	return func(ctx context.Context, data NewStatisticData) (Statistic, error) {
		if err := data.validate(); err != nil {
			return Statistic{}, err
		}

		return storageCreateStatisticFunc(ctx, data)
	}
}

func BuildGetStatisticByIDAndGameID(storageGetStatisticByIDAndGameID StorageGetStatisticByIDAndGameID) GetByIDAndGameIDFunc {
	return func(ctx context.Context, id, gameID string) (Statistic, error) {
		return storageGetStatisticByIDAndGameID(ctx, id, gameID)
	}
}

func BuildSoftDeleteStatistic(storageSoftDeleteStatistic StorageSoftDeleteStatistic) SoftDeleteByIDAndGameIDFunc {
	return func(ctx context.Context, id, gameID string) error {
		return storageSoftDeleteStatistic(ctx, id, gameID)
	}
}
