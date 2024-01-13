package rest

import (
	"time"

	"github.com/gabarcia/metagaming-api/internal/quest"
)

type Task struct {
	CreatedAt   time.Time `json:"createdAt"`           // Time that the task was created
	UpdatedAt   time.Time `json:"updatedAt"`           // Last time that the task was updated
	ID          string    `json:"id"`                  // Task ID
	Name        string    `json:"name"`                // Task name
	Description string    `json:"description"`         // Task details
	DependsOn   string    `json:"dependsOn,omitempty"` // ID of the task that needs to be completed before this one can be started
	Rule        string    `json:"rule"`                // Task completion logic as JsonLogic. See https://jsonlogic.com/
}

func taskFromDomain(t quest.Task) Task {
	return Task{
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		DependsOn:   t.DependsOn,
		Rule:        t.Rule,
	}
}
