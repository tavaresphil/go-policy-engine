package native_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tavaresphil/go-policy-engine/pkg/evaluators/native"
	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

func TestTemporalHandler_Eval(t *testing.T) {
	type input struct {
		pc   policies.PolicyCondition
		attr policies.Resolver
	}
	type output struct {
		res bool
		err error
	}

	before := "2020-01-01"
	after := "2020-01-02"
	unixBefore := 1577836800 // 2020-01-01
	unixAfter := 1577923200  // 2020-01-02
	tm := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		input  input
		output output
		assert func(t *testing.T, expected, actual output)
	}{
		{
			name: "when attribute not found should return error",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "ts", Operator: policies.OpBefore, Value: before},
				attr: policies.MapAttributes{},
			},
			output: output{res: false, err: fmt.Errorf("missing required attribute: ts")},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "when attribute before string should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "ts", Operator: policies.OpBefore, Value: after},
				attr: policies.MapAttributes{"ts": before},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when attribute after string should return true for OpAfter",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "ts", Operator: policies.OpAfter, Value: before},
				attr: policies.MapAttributes{"ts": after},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when attribute unix ints should compare as time",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "ts", Operator: policies.OpAfter, Value: unixBefore},
				attr: policies.MapAttributes{"ts": unixAfter},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when attribute time.Time should compare",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "ts", Operator: policies.OpBefore, Value: tm},
				attr: policies.MapAttributes{"ts": before},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when value cannot be parsed should return error",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "ts", Operator: policies.OpBefore, Value: "not a date"},
				attr: policies.MapAttributes{"ts": before},
			},
			output: output{res: false, err: fmt.Errorf("unable to parse time string: not a date")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "unable to parse time string")
				}
			},
		},
		{
			name: "when unsupported operator should return error",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "ts", Operator: "unknown", Value: after},
				attr: policies.MapAttributes{"ts": before},
			},
			output: output{res: false, err: fmt.Errorf("unsupported temporal operator: %s", "unknown")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "unsupported temporal operator")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdl := native.NewTemporalHandler()

			var actual output
			actual.res, actual.err = hdl.Eval(tt.input.pc, tt.input.attr)

			tt.assert(t, tt.output, actual)
		})
	}
}
