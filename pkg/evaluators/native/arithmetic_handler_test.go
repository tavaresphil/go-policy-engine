package native_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tavaresphil/go-policy-engine/pkg/evaluators/native"
	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

func TestArithmeticHandler_Eval(t *testing.T) {
	type input struct {
		pc   policies.PolicyCondition
		attr policies.Resolver
	}
	type output struct {
		res bool
		err error
	}

	tests := []struct {
		name   string
		input  input
		output output
		assert func(t *testing.T, expected, actual output)
	}{
		{
			name: "when attribute not found should return error",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpMod, Value: 2},
				attr: policies.MapAttributes{},
			},
			output: output{res: false, err: fmt.Errorf("missing required attribute: n")},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "when divisible should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpMod, Value: 2},
				attr: policies.MapAttributes{"n": 10},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when not divisible should return false",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpMod, Value: 3},
				attr: policies.MapAttributes{"n": 10},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when float modulus works",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpMod, Value: 0.5},
				attr: policies.MapAttributes{"n": 10.5},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when divisor zero should return error",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpMod, Value: 0},
				attr: policies.MapAttributes{"n": 10},
			},
			output: output{res: false, err: fmt.Errorf("modulo by zero")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "modulo by zero")
				}
			},
		},
		{
			name: "when numeric strings should work",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "n", Operator: policies.OpMod, Value: "2"},
				attr: policies.MapAttributes{"n": "10"},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdl := native.NewArithmeticHandler()

			var actual output
			actual.res, actual.err = hdl.Eval(tt.input.pc, tt.input.attr)

			tt.assert(t, tt.output, actual)
		})
	}
}
