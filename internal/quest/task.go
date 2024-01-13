package quest

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/diegoholiveira/jsonlogic/v3"
)

var (
	ErrTaskValidationError          = errors.New("task validation error")
	ErrInvalidTaskName              = errors.New("invalid task name")
	ErrInvalidTaskRule              = errors.New("invalid task rule")
	ErrTaskRuleNotBoolean           = errors.New("task rule does not return a boolean value")
	ErrBrokenTaskRuleData           = errors.New("broken rule data")
	ErrInvalidSucessRuleDataExemple = errors.New("success exemple task rule data returned false")
)

type Task struct {
	ID          string // Task ID
	Name        string // Task name
	Description string // Task details
	DependsOn   string // ID of the task that needs to be completed before this one can be started
	Rule        string // Task completion logic as JsonLogic. See https://jsonlogic.com/
}

func (t Task) validate(successExempleData string) error {
	errList := make([]error, 0)

	if t.Name == "" {
		errList = append(errList, ErrInvalidTaskName)
	}

	if !jsonlogic.IsValid(strings.NewReader(t.Rule)) {
		errList = append(errList, ErrInvalidTaskRule)
	}

	ok, err := t.apply(successExempleData)
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

func (t Task) apply(v string) (bool, error) {
	var (
		rule = strings.NewReader(t.Rule)
		data = strings.NewReader(v)
	)

	var r bytes.Buffer
	if err := jsonlogic.Apply(rule, data, &r); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			err = ErrBrokenTaskRuleData
		}

		return false, err
	}

	boolValue, err := strconv.ParseBool(strings.TrimSpace(r.String()))
	if err != nil {
		return false, ErrTaskRuleNotBoolean
	}

	return boolValue, nil
}
