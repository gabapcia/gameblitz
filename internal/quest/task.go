package quest

import (
	"errors"
	"time"
)

var (
	ErrTaskValidationError          = errors.New("task validation error")
	ErrInvalidTaskName              = errors.New("invalid task name")
	ErrInvalidTaskRule              = errors.New("invalid task rule")
	ErrInvalidSucessRuleDataExemple = errors.New("success exemple task rule data returned false")
)

type NewTaskData struct {
	Name        string // Task name
	Description string // Task details
	DependsOn   *int   // Array index of the task that needs to be completed before this one can be started
	Rule        string // Task completion logic as JsonLogic. See https://jsonlogic.com/
}

type Task struct {
	CreatedAt   time.Time // Time that the task was created
	UpdatedAt   time.Time // Last time that the task was updated
	DeletedAt   time.Time // Time that the task was deleted
	ID          string    // Task ID
	Name        string    // Task name
	Description string    // Task details
	DependsOn   string    // ID of the task that needs to be completed before this one can be started
	Rule        string    // Task completion logic as JsonLogic. See https://jsonlogic.com/
}

func (t NewTaskData) validate(successExempleData string) error {
	errList := make([]error, 0)

	if t.Name == "" {
		errList = append(errList, ErrInvalidTaskName)
	}

	if !RuleIsValid(t.Rule) {
		errList = append(errList, ErrInvalidTaskRule)
	}

	ok, err := RuleApply(t.Rule, successExempleData)
	if err != nil {
		errList = append(errList, err)
	}

	if !ok {
		errList = append(errList, ErrInvalidSucessRuleDataExemple)
	}

	if len(errList) > 0 {
		errList = append(errList, ErrTaskValidationError)
	}

	return errors.Join(errList...)
}
