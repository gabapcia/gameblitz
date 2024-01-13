package quest

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrQuestValidationError               = errors.New("invalid quest")
	ErrInvalidQuestName                   = errors.New("invalid quest name")
	ErrQuestMissingGameID                 = errors.New("missing game id")
	ErrQuestTaskRuleSuceessDataIncomplete = errors.New("missing success data for some tasks")
)

type NewQuestData struct {
	GameID          string        // ID of the game responsible for the quest
	Name            string        // Quest name
	Description     string        // Quest details
	Tasks           []NewTaskData // Quest task list
	TasksValidators []string      // Quest task list success validation data
}

type Quest struct {
	CreatedAt   time.Time // Time that the quest was created
	UpdatedAt   time.Time // Last time that the quest was updated
	DeletedAt   time.Time // Time that the quest was deleted
	ID          string    // Quest ID
	GameID      string    // ID of the game responsible for the quest
	Name        string    // Quest name
	Description string    // Quest details
	Tasks       []Task    // Quest task list
}

func (q NewQuestData) validate() error {
	errList := make([]error, 0)

	if q.Name == "" {
		errList = append(errList, ErrInvalidQuestName)
	}

	if q.GameID == "" {
		errList = append(errList, ErrQuestMissingGameID)
	}

	if len(q.Tasks) != len(q.TasksValidators) {
		errList = append(errList, ErrQuestTaskRuleSuceessDataIncomplete)
	} else {
		for i, task := range q.Tasks {
			if err := task.validate(q.TasksValidators[i]); err != nil {
				errList = append(errList, fmt.Errorf("Task #%d\n%w", i, err))
			}
		}
	}

	if len(errList) > 0 {
		errList = append(errList, ErrQuestValidationError)
	}

	return errors.Join(errList...)
}

func BuildCreateQuestFunc(storageCreateQuestFunc StorageCreateQuestFunc) CreateQuestFunc {
	return func(ctx context.Context, data NewQuestData) (Quest, error) {
		if err := data.validate(); err != nil {
			return Quest{}, err
		}

		return storageCreateQuestFunc(ctx, data)
	}
}
