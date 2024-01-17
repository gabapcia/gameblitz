package quest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskDepencyIsCyclic(t *testing.T) {
	t.Run("Empty Slice", func(t *testing.T) {
		isCyclic := taskDepencyIsCyclic([]NewTaskData{})
		assert.False(t, isCyclic)
	})

	t.Run("No Cycle", func(t *testing.T) {
		// A <-- B <-- C <-- +
		//       |           |
		//       + <-- D <-- E
		isCyclic := taskDepencyIsCyclic([]NewTaskData{
			{Name: "A", DependsOn: []int{}},
			{Name: "B", DependsOn: []int{0}},
			{Name: "C", DependsOn: []int{1}},
			{Name: "D", DependsOn: []int{1}},
			{Name: "E", DependsOn: []int{2, 3}},
		})

		assert.False(t, isCyclic)
	})

	t.Run("Self Dependency", func(t *testing.T) {
		isCyclic := taskDepencyIsCyclic([]NewTaskData{
			{Name: "A", DependsOn: []int{0}},
		})

		assert.True(t, isCyclic)
	})

	t.Run("Simple Cycle", func(t *testing.T) {
		// A <-- B <-- C
		// + --------> +
		isCyclic := taskDepencyIsCyclic([]NewTaskData{
			{Name: "A", DependsOn: []int{2}},
			{Name: "B", DependsOn: []int{0}},
			{Name: "C", DependsOn: []int{1}},
		})

		assert.True(t, isCyclic)
	})

	t.Run("Complex Cycle", func(t *testing.T) {
		// A <-- B <-- C <-- +
		// |     |     	     |
		// |     + <-- D --- +
		// |           |
		// |           + <-- E
		// |                 |
		// + --------------> +
		isCyclic := taskDepencyIsCyclic([]NewTaskData{
			{Name: "A", DependsOn: []int{4}},
			{Name: "B", DependsOn: []int{0}},
			{Name: "C", DependsOn: []int{1}},
			{Name: "D", DependsOn: []int{1, 2}},
			{Name: "E", DependsOn: []int{3}},
		})

		assert.True(t, isCyclic)
	})

	t.Run("Ping Pong", func(t *testing.T) {
		// A <-- B
		// |     |
		// + --> +
		isCyclic := taskDepencyIsCyclic([]NewTaskData{
			{Name: "A", DependsOn: []int{1}},
			{Name: "B", DependsOn: []int{0}},
		})

		assert.True(t, isCyclic)
	})
}

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
