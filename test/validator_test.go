package test

import (
	"fmt"
	"github.com/lanvard/errors"
	"github.com/lanvard/support"
	"github.com/lanvard/syslog/log_level"
	"github.com/lanvard/validation/rule"
	"github.com/lanvard/validation/val"
	"github.com/stretchr/testify/require"
	net "net/http"
	"testing"
)

func Test_validate_nothing(t *testing.T) {
	errs := val.Validate(support.NewValue(nil))
	require.Equal(t, []error{}, errs)
}

func Test_validate_nothing_with_empty_verification(t *testing.T) {
	errs := val.Validate(
		support.NewValue(nil),
		val.Verify("title"),
	)
	require.Empty(t, errs)
}

func Test_validate_nothing_with_empty_verifications(t *testing.T) {
	errs := val.Validate(
		support.NewValue(nil),
		val.Verify("title"),
		val.Verify("description"),
	)
	require.Empty(t, errs)
}

func Test_validate_with_multiple_verifications(t *testing.T) {
	errs := val.Validate(
		support.NewValue(map[string]string{"title": "Horse", "description": "Big animal"}),
		val.Verify("title"),
		val.Verify("description"),
	)
	require.Empty(t, errs)
}

func Test_validate_with_multiple_invalid_keys(t *testing.T) {
	errs := val.Validate(
		support.NewValue(map[string]string{}),
		val.Verify("title", rule.Required{}),
		val.Verify("description", rule.Required{}),
	)
	require.Len(t, errs, 2)
}

func Test_validate_invalid_values_with_multiple_rules(t *testing.T) {
	errs := val.Validate(
		support.NewValue(nil),
		val.Verify("title", rule.Present{}, rule.Required{}),
	)
	require.Len(t, errs, 1)
	require.EqualError(t, errs[0], "field title must be present")
}

func Test_validate_nested_key_error(t *testing.T) {
	errs := val.Validate(
		support.NewValue(map[string]string{}),
		val.Verify("user.title", rule.Present{}),
	)
	require.EqualError(t, errs[0], "field user.title must be present")
}

func Test_validate_map(t *testing.T) {
	errs := val.Validate(
		map[string]string{},
		val.Verify("user.title", rule.Present{}),
	)
	require.EqualError(t, errs[0], "field user.title must be present")
}

func Test_error_has_stack_trace(t *testing.T) {
	errs := val.Validate(
		map[string]string{},
		val.Verify("user.title", rule.Present{}),
	)
	stack, ok := errors.FindStack(errs[0])
	require.True(t, ok)
	require.Contains(t, fmt.Sprintf("%+v", stack), "validator_test.go")
}

func Test_normal_rule_not_required(t *testing.T) {
	errs := val.Validate(
		nil,
		val.Verify("title", mockRuleNotRequired{}),
	)
	require.Empty(t, errs)
}

func Test_validation_error_status(t *testing.T) {
	errs := val.Validate(
		map[string]string{},
		val.Verify("user.title", rule.Present{}),
	)
	status, _ := errors.FindStatus(errs[0])
	require.Equal(t, net.StatusUnprocessableEntity, status)
}

func Test_validation_log_level(t *testing.T) {
	errs := val.Validate(
		map[string]string{},
		val.Verify("user.title", rule.Present{}),
	)
	level, _ := errors.FindLevel(errs[0])
	require.Equal(t, log_level.INFO, level)
}

type mockRuleNotRequired struct{}

func (m mockRuleNotRequired) Verify(value support.Value) error {
	return errors.New("don't show this error if value not present")
}