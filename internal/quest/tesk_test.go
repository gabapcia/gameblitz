package quest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskValidate(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		data := `{"killed": {"terrorists": 200}}`
		task := NewTaskData{
			Name:        "Test Task",
			Description: "Test task description",
			DependsOn:   nil,
			Rule:        `{">": [{"var": "killed.terrorists"}, 150]}`,
		}

		err := task.validate(data)
		assert.NoError(t, err)
	})

	t.Run("Validation Error", func(t *testing.T) {
		data := `{}`
		task := NewTaskData{}

		err := task.validate(data)
		assert.ErrorIs(t, err, ErrTaskValidationError)
		assert.ErrorIs(t, err, ErrInvalidTaskName)
		assert.ErrorIs(t, err, ErrInvalidTaskRule)
		assert.ErrorIs(t, err, ErrInvalidSucessRuleDataExemple)
		assert.ErrorIs(t, err, ErrBrokenRuleData)
	})
}
