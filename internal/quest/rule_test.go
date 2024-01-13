package quest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuleIsValid(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		rule := `{">": [{"var": "killed.terrorists"}, 150]}`

		isValid := RuleIsValid(rule)
		assert.True(t, isValid)
	})

	t.Run("Broken", func(t *testing.T) {
		rule := `{`

		isValid := RuleIsValid(rule)
		assert.False(t, isValid)
	})
}

func TestRuleApply(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		var (
			rule = `{">": [{"var": "killed.terrorists"}, 150]}`
			data = `{"killed": {"terrorists": 200}}`
		)

		boolValue, err := RuleApply(rule, data)
		assert.NoError(t, err)
		assert.True(t, boolValue)
	})

	t.Run("Broken Data", func(t *testing.T) {
		var (
			rule = `{">": [{"var": "killed.terrorists"}, 150]}`
			data = `{`
		)

		boolValue, err := RuleApply(rule, data)
		assert.ErrorIs(t, err, ErrBrokenRuleData)
		assert.False(t, boolValue)
	})

	t.Run("Rule Not Return Boolean Value", func(t *testing.T) {
		var (
			rule = `{
				"filter": [
					{"var": "users"},
					{">=": [{"var": ".age"}, 18]}
				]
			}`
			data = `{
				"users": [
					{"name": "Diego", "age": 33, "location": "Florianópolis"},
					{"name": "Jack", "age": 12, "location": "London"},
					{"name": "Pedro", "age": 19, "location": "Lisbon"},
					{"name": "Leopoldina", "age": 30, "location": "Rio de Janeiro"}
				]
			}`
		)

		boolValue, err := RuleApply(rule, data)
		assert.ErrorIs(t, err, ErrRuleNotBoolean)
		assert.False(t, boolValue)
	})
}
