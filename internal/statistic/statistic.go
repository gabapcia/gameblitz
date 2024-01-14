package statistic

import (
	"context"
	"errors"
	"slices"
	"time"
)

var (
	ErrStatisticValidation            = errors.New("invalid statistic")
	ErrInvalidName                    = errors.New("invalid name")
	ErrMissingGameID                  = errors.New("missing game id")
	ErrInvalidAggregationMode         = errors.New("invalid aggregation mode")
	ErrInvalidLandmarkLowerThanGoal   = errors.New("landmark lower than goal")
	ErrInvalidLandmarkGreaterThanGoal = errors.New("landmark greater than goal")
	ErrCannotOverflowWithNoGoal       = errors.New("cannot overflow but there's no goal")
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
	CanOverflow     bool      // Can overflow the goal?
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
	CanOverflow     bool      // Can overflow the goal?
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

	if s.Goal == nil && !s.CanOverflow {
		errList = append(errList, ErrCannotOverflowWithNoGoal)
	}

	if !slices.Contains(AggregationModes, s.AggregationMode) {
		errList = append(errList, ErrInvalidAggregationMode)
	} else if s.Goal != nil && !s.CanOverflow {
		for _, landmark := range s.Landmarks {
			switch s.AggregationMode {
			case AggregationModeMin, AggregationModeSub:
				if landmark < *s.Goal {
					errList = append(errList, ErrInvalidLandmarkLowerThanGoal)
					break
				}
			case AggregationModeMax, AggregationModeSum:
				if landmark > *s.Goal {
					errList = append(errList, ErrInvalidLandmarkGreaterThanGoal)
					break
				}
			}
		}
	}

	if len(errList) > 0 {
		errList = append(errList, ErrStatisticValidation)
	}

	return errors.Join(errList...)
}

func BuildCreateStatisticFunc(storageCreateStatisticFunc StorageCreateStatisticFunc) CreateStatisticFunc {
	return func(ctx context.Context, data NewStatisticData) (Statistic, error) {
		if err := data.validate(); err != nil {
			return Statistic{}, err
		}

		return storageCreateStatisticFunc(ctx, data)
	}
}

func BuildGetStatisticByIDAndGameID(storageGetStatisticByIDAndGameID StorageGetStatisticByIDAndGameID) GetStatisticByIDAndGameID {
	return func(ctx context.Context, id, gameID string) (Statistic, error) {
		return storageGetStatisticByIDAndGameID(ctx, id, gameID)
	}
}

func BuildSoftDeleteStatistic(storageSoftDeleteStatistic StorageSoftDeleteStatistic) SoftDeleteStatistic {
	return func(ctx context.Context, id, gameID string) error {
		return storageSoftDeleteStatistic(ctx, id, gameID)
	}
}
