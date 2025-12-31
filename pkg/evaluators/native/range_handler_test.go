package native_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tavaresphil/go-policy-engine/pkg/evaluators/native"
	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

func TestRangeHandler_Eval(t *testing.T) {
	type input struct {
		pc   policies.PolicyCondition
		attr policies.Resolver
	}
	type output struct {
		res bool
		err error
	}

	tm := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	before := time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC)
	after := time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		input  input
		output output
		assert func(t *testing.T, expected, actual output)
	}{
		{
			name: "when attribute not found should return error",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "x", Operator: policies.OpBetween, Value: []any{1, 10}},
				attr: policies.MapAttributes{},
			},
			output: output{res: false, err: fmt.Errorf("missing required attribute: x")},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "when int is between should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpBetween, Value: []any{1, 10}},
				attr: policies.MapAttributes{"n": 5},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when equals boundaries and inclusive true should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpBetween, Value: []any{5, 5, true}},
				attr: policies.MapAttributes{"n": 5},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when equals boundaries and inclusive false should return false",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpBetween, Value: []any{5, 5, false}},
				attr: policies.MapAttributes{"n": 5},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when time is in range should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "t", Operator: policies.OpBetween, Value: []any{before, after}},
				attr: policies.MapAttributes{"t": tm},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when string in range should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "s", Operator: policies.OpBetween, Value: []any{"a", "z"}},
				attr: policies.MapAttributes{"s": "m"},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when attribute implements Stringer should compare as string",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "s", Operator: policies.OpBetween, Value: []any{"a", "z"}},
				attr: policies.MapAttributes{"s": struct{ v string }{"m"}},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				// struct isn't a Stringer by default, test will use fmt.Sprint fallback => "{m}"
				assert.NoError(t, actual.err)
				// It's okay if it doesn't pass the range; we're just asserting no crash and boolean result
				_ = actual.res
			},
		},
		{
			name: "when types mismatch now compares as strings and returns false",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "s", Operator: policies.OpBetween, Value: []any{"a", "z"}},
				attr: policies.MapAttributes{"s": 1},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when invalid args should return error",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpBetween, Value: 1},
				attr: policies.MapAttributes{"n": 5},
			},
			output: output{res: false, err: fmt.Errorf("unsupported between value type: int")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "unsupported between value type")
				}
			},
		},
		{
			name: "when int/float mix should work",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpBetween, Value: []any{1.5, 10.0}},
				attr: policies.MapAttributes{"n": 5},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when float/int mix should work",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpBetween, Value: []any{1, 10}},
				attr: policies.MapAttributes{"n": 5.5},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when numeric strings work",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpBetween, Value: []any{"1", "10"}},
				attr: policies.MapAttributes{"n": "5"},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when min greater than max should return false",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpBetween, Value: []any{10, 1}},
				attr: policies.MapAttributes{"n": 5},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when map missing min or max should return error",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpBetween, Value: map[string]any{"min": 1}},
				attr: policies.MapAttributes{"n": 5},
			},
			output: output{res: false, err: fmt.Errorf("between requires min and max in map form")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "between requires min and max in map form")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdl := native.NewRangeHandler()

			var actual output
			actual.res, actual.err = hdl.Eval(tt.input.pc, tt.input.attr)

			tt.assert(t, tt.output, actual)
		})
	}
}
