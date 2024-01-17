package quest

import (
	"errors"
	"slices"
	"time"
)

var (
	ErrTaskValidationError          = errors.New("task validation error")
	ErrInvalidTaskID                = errors.New("invalid task id")
	ErrInvalidTaskName              = errors.New("invalid task name")
	ErrInvalidTaskRule              = errors.New("invalid task rule")
	ErrInvalidSucessRuleDataExemple = errors.New("success exemple task rule data returned false")
	ErrInvalidTaskDependencyIndex   = errors.New("invalid task dependency array index")
	ErrTaskDependencyCycle          = errors.New("task dependency cycle detected")
)

type NewTaskData struct {
	Name                  string // Task name
	Description           string // Task details
	DependsOn             []int  // List of array indexes of the tasks that needs to be completed before this one can be started
	RequiredForCompletion bool   // Is this task required for the quest completion?
	Rule                  string // Task completion logic as JsonLogic. See https://jsonlogic.com/
}

type Task struct {
	CreatedAt             time.Time // Time that the task was created
	UpdatedAt             time.Time // Last time that the task was updated
	DeletedAt             time.Time // Time that the task was deleted
	ID                    string    // Task ID
	Name                  string    // Task name
	Description           string    // Task details
	DependsOn             []string  // IDs from the tasks that needs to be completed before this one can be started
	RequiredForCompletion bool      // Is this task required for the quest completion?
	Rule                  string    // Task completion logic as JsonLogic. See https://jsonlogic.com/
}

func taskDepencyHasCycle(tasks []NewTaskData, index int, visited map[int]bool, recStack map[int]bool) bool {
	visited[index] = true
	recStack[index] = true

	for _, depIndex := range tasks[index].DependsOn {
		if !visited[depIndex] {
			if taskDepencyHasCycle(tasks, depIndex, visited, recStack) {
				return true
			}
		} else if recStack[depIndex] {
			return true
		}
	}

	recStack[index] = false
	return false
}

func taskDepencyIsCyclic(tasks []NewTaskData) bool {
	var (
		visited  = make(map[int]bool)
		recStack = make(map[int]bool)
	)

	for i := range tasks {
		if !visited[i] {
			if taskDepencyHasCycle(tasks, i, visited, recStack) {
				return true
			}
		}
	}

	return false
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
		errList = slices.Insert(errList, 0, ErrTaskValidationError)
	}

	return errors.Join(errList...)
}
