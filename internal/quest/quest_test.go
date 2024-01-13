package quest

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestQuestValidate(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		quest := NewQuestData{
			GameID:      uuid.NewString(),
			Name:        "Test Quest",
			Description: "Quest description",
			Tasks: []NewTaskData{
				{
					Name:        "Test Task",
					Description: "Test task description",
					DependsOn:   nil,
					Rule:        `{">": [{"var": "killed.terrorists"}, 150]}`,
				},
			},
			TasksValidators: []string{
				`{"killed": {"terrorists": 200}}`,
			},
		}

		err := quest.validate()
		assert.NoError(t, err)
	})

	t.Run("Validation Error", func(t *testing.T) {
		quest := NewQuestData{
			Tasks: []NewTaskData{
				{},
			},
			TasksValidators: []string{
				`{}`,
			},
		}

		err := quest.validate()
		// Quest errors
		assert.ErrorIs(t, err, ErrQuestValidationError)
		assert.ErrorIs(t, err, ErrQuestMissingGameID)
		assert.ErrorIs(t, err, ErrInvalidQuestName)
		// Task errors
		assert.ErrorIs(t, err, ErrTaskValidationError)
		assert.ErrorIs(t, err, ErrInvalidTaskName)
		assert.ErrorIs(t, err, ErrInvalidTaskRule)
		assert.ErrorIs(t, err, ErrInvalidSucessRuleDataExemple)
		assert.ErrorIs(t, err, ErrBrokenRuleData)
	})

	t.Run("Missing Task Success Data Exemple", func(t *testing.T) {
		quest := NewQuestData{
			GameID:      uuid.NewString(),
			Name:        "Test Quest",
			Description: "Quest description",
			Tasks: []NewTaskData{
				{
					Name:        "Test Task",
					Description: "Test task description",
					DependsOn:   nil,
					Rule:        `{">": [{"var": "killed.terrorists"}, 150]}`,
				},
			},
			TasksValidators: []string{},
		}

		err := quest.validate()
		// Quest errors
		assert.ErrorIs(t, err, ErrQuestValidationError)
		assert.ErrorIs(t, err, ErrQuestTaskRuleSuceessDataIncomplete)
		// Task errors
		assert.NotErrorIs(t, err, ErrTaskValidationError)
	})
}
