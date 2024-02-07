package rest

import (
	"time"

	"github.com/gabapcia/gameblitz/internal/quest"
)

type Task struct {
	CreatedAt             time.Time `json:"createdAt"`             // Time that the task was created
	UpdatedAt             time.Time `json:"updatedAt"`             // Last time that the task was updated
	ID                    string    `json:"id"`                    // Task ID
	Name                  string    `json:"name"`                  // Task name
	Description           string    `json:"description"`           // Task details
	DependsOn             []string  `json:"dependsOn,omitempty"`   // IDs from the tasks that needs to be completed before this one can be started
	RequiredForCompletion bool      `json:"requiredForCompletion"` // Is this task required for the quest completion?
	Rule                  string    `json:"rule"`                  // Task completion logic as JsonLogic. See https://jsonlogic.com/
}

func taskFromDomain(t quest.Task) Task {
	return Task{
		CreatedAt:             t.CreatedAt,
		UpdatedAt:             t.UpdatedAt,
		ID:                    t.ID,
		Name:                  t.Name,
		Description:           t.Description,
		DependsOn:             t.DependsOn,
		RequiredForCompletion: t.RequiredForCompletion,
		Rule:                  t.Rule,
	}
}
