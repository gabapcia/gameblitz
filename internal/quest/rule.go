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
	ErrRuleNotBoolean = errors.New("rule does not return a boolean value")
	ErrBrokenRuleData = errors.New("broken rule data")
)

func RuleIsValid(r string) bool {
	return jsonlogic.IsValid(strings.NewReader(r))
}

func RuleApply(r, v string) (bool, error) {
	var (
		rule = strings.NewReader(r)
		data = strings.NewReader(v)
	)

	var b bytes.Buffer
	if err := jsonlogic.Apply(rule, data, &b); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			err = ErrBrokenRuleData
		}

		return false, err
	}

	boolValue, err := strconv.ParseBool(strings.TrimSpace(b.String()))
	if err != nil {
		return false, ErrRuleNotBoolean
	}

	return boolValue, nil
}
