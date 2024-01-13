package quest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskApply(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		task := Task{
			Name:        "Test Task",
			Description: "Test task description",
			DependsOn:   "",
			Rule:        `{">": [{"var": "killed.terrorists"}, 150]}`,
		}

		data := `{"killed": {"terrorists": 200}}`

		boolValue, err := task.apply(data)
		assert.NoError(t, err)
		assert.True(t, boolValue)
	})

	t.Run("Broken Data", func(t *testing.T) {
		task := Task{
			Name:        "Test Task",
			Description: "Test task description",
			DependsOn:   "",
			Rule:        `{">": [{"var": "killed.terrorists"}, 150]}`,
		}

		data := `{`

		boolValue, err := task.apply(data)
		assert.ErrorIs(t, err, ErrBrokenTaskRuleData)
		assert.False(t, boolValue)
	})

	t.Run("Rule Not Return Boolean Value", func(t *testing.T) {
		task := Task{
			Name:        "Test Task",
			Description: "Test task description",
			DependsOn:   "",
			Rule: `{
				"filter": [
					{"var": "users"},
					{">=": [{"var": ".age"}, 18]}
				]
			}`,
		}

		data := `{
			"users": [
				{"name": "Diego", "age": 33, "location": "Florianópolis"},
				{"name": "Jack", "age": 12, "location": "London"},
				{"name": "Pedro", "age": 19, "location": "Lisbon"},
				{"name": "Leopoldina", "age": 30, "location": "Rio de Janeiro"}
			]
		}`

		boolValue, err := task.apply(data)
		assert.ErrorIs(t, err, ErrTaskRuleNotBoolean)
		assert.False(t, boolValue)
	})
}

func TestTaskValidate(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		data := `{"killed": {"terrorists": 200}}`
		task := Task{
			Name:        "Test Task",
			Description: "Test task description",
			DependsOn:   "",
			Rule:        `{">": [{"var": "killed.terrorists"}, 150]}`,
		}

		err := task.validate(data)
		assert.NoError(t, err)
	})

	t.Run("Validation Error", func(t *testing.T) {
		data := `{}`
		task := Task{}

		err := task.validate(data)
		assert.ErrorIs(t, err, ErrTaskValidationError)
		assert.ErrorIs(t, err, ErrInvalidTaskName)
		assert.ErrorIs(t, err, ErrInvalidTaskRule)
		assert.ErrorIs(t, err, ErrInvalidSucessRuleDataExemple)
		assert.ErrorIs(t, err, ErrBrokenTaskRuleData)
	})
}
