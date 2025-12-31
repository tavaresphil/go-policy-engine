package native_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tavaresphil/go-policy-engine/pkg/evaluators/native"
	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

func TestNativeEngine_Eval(t *testing.T) {
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
			name:   "when operator unknown should return error",
			input:  input{pc: policies.PolicyCondition{Operator: policies.Operator("foo")}},
			output: output{res: false, err: fmt.Errorf("unknown operator: foo")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "unknown operator")
				}
			},
		},
		{
			name: "when attribute not found should return error",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "has_dependency", Operator: policies.OpEqual, Value: true},
				attr: policies.MapAttributes{},
			},
			output: output{res: false, err: fmt.Errorf("missing required attribute: has_dependency")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "missing required attribute")
				}
			},
		},
		{
			name:   "when comparison equals should return true",
			input:  input{pc: policies.PolicyCondition{Attribute: "active", Operator: policies.OpEqual, Value: true}, attr: policies.MapAttributes{"active": true}},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name:   "when mod operator satisfied should return true",
			input:  input{pc: policies.PolicyCondition{Attribute: "count", Operator: policies.OpMod, Value: 3}, attr: policies.MapAttributes{"count": 9}},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name:   "when string contains should return true",
			input:  input{pc: policies.PolicyCondition{Attribute: "name", Operator: policies.OpContains, Value: "oh"}, attr: policies.MapAttributes{"name": "oh hai"}},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name:   "when element is in set should return true",
			input:  input{pc: policies.PolicyCondition{Attribute: "elem", Operator: policies.OpIn, Value: []int{1, 2, 3}}, attr: policies.MapAttributes{"elem": 2}},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name:   "when number is between should return true",
			input:  input{pc: policies.PolicyCondition{Attribute: "n", Operator: policies.OpBetween, Value: []int{1, 3}}, attr: policies.MapAttributes{"n": 2}},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name:   "when temporal before should return true",
			input:  input{pc: policies.PolicyCondition{Attribute: "ts", Operator: policies.OpBefore, Value: "2021-01-01"}, attr: policies.MapAttributes{"ts": "2020-01-01"}},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name:   "when logical and should combine children",
			input:  input{pc: policies.PolicyCondition{Operator: policies.OpAnd, Conditions: []policies.PolicyCondition{{Attribute: "a", Operator: policies.OpEqual, Value: 1}, {Attribute: "b", Operator: policies.OpEqual, Value: 2}}}, attr: policies.MapAttributes{"a": 1, "b": 2}},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eng := native.NewNativeEngine()

			var actual output
			actual.res, actual.err = eng.Eval(tt.input.pc, tt.input.attr)

			tt.assert(t, tt.output, actual)
		})
	}
}
