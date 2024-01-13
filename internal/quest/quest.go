package quest

import (
	"errors"
	"fmt"
)

var (
	ErrQuestValidationError               = errors.New("invalid quest")
	ErrInvalidQuestName                   = errors.New("invalid quest name")
	ErrQuestMissingGameID                 = errors.New("missing game id")
	ErrQuestTaskRuleSuceessDataIncomplete = errors.New("missing success data for some tasks")
)

type Quest struct {
	ID          string // Quest ID
	GameID      string // ID of the game responsible for the quest
	Name        string // Quest name
	Description string // Quest details
	Tasks       []Task // Quest task list
}

func (q Quest) validate(successExempleData []string) error {
	errList := make([]error, 0)

	if q.Name == "" {
		errList = append(errList, ErrInvalidQuestName)
	}

	if q.GameID == "" {
		errList = append(errList, ErrQuestMissingGameID)
	}

	if len(q.Tasks) != len(successExempleData) {
		errList = append(errList, ErrQuestTaskRuleSuceessDataIncomplete)
	} else {
		for i, task := range q.Tasks {
			if err := task.validate(successExempleData[i]); err != nil {
				errList = append(errList, fmt.Errorf("Task #%d: %w", i, err))
			}
		}
	}

	if len(errList) > 0 {
		errList = append(errList, ErrQuestValidationError)
	}

	return errors.Join(errList...)
}
