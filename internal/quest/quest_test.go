package quest

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestQuestValidate(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		data := []string{
			`{"killed": {"terrorists": 200}}`,
		}
		quest := Quest{
			GameID:      uuid.NewString(),
			Name:        "Test Quest",
			Description: "Quest description",
			Tasks: []Task{
				{
					Name:        "Test Task",
					Description: "Test task description",
					DependsOn:   "",
					Rule:        `{">": [{"var": "killed.terrorists"}, 150]}`,
				},
			},
		}

		err := quest.validate(data)
		assert.NoError(t, err)
	})

	t.Run("Validation Error", func(t *testing.T) {
		data := []string{
			`{}`,
		}
		quest := Quest{
			Tasks: []Task{
				{},
			},
		}

		err := quest.validate(data)
		// Quest errors
		assert.ErrorIs(t, err, ErrQuestValidationError)
		assert.ErrorIs(t, err, ErrQuestMissingGameID)
		assert.ErrorIs(t, err, ErrInvalidQuestName)
		// Task errors
		assert.ErrorIs(t, err, ErrTaskValidationError)
		assert.ErrorIs(t, err, ErrInvalidTaskName)
		assert.ErrorIs(t, err, ErrInvalidTaskRule)
		assert.ErrorIs(t, err, ErrInvalidSucessRuleDataExemple)
		assert.ErrorIs(t, err, ErrBrokenTaskRuleData)
	})

	t.Run("Missing Task Success Data Exemple", func(t *testing.T) {
		data := []string{}
		quest := Quest{
			GameID:      uuid.NewString(),
			Name:        "Test Quest",
			Description: "Quest description",
			Tasks: []Task{
				{
					Name:        "Test Task",
					Description: "Test task description",
					DependsOn:   "",
					Rule:        `{">": [{"var": "killed.terrorists"}, 150]}`,
				},
			},
		}

		err := quest.validate(data)
		// Quest errors
		assert.ErrorIs(t, err, ErrQuestValidationError)
		assert.ErrorIs(t, err, ErrQuestTaskRuleSuceessDataIncomplete)
		// Task errors
		assert.NotErrorIs(t, err, ErrTaskValidationError)
	})
}
