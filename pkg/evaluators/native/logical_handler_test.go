package native_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tavaresphil/go-policy-engine/pkg/evaluators/native"
	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

func TestLogicalHandler_Eval(t *testing.T) {
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
			name: "when all children true should return true",
			input: input{
				pc: policies.PolicyCondition{
					Operator: policies.OpAnd,
					Conditions: []policies.PolicyCondition{
						{Attribute: "a", Operator: policies.OpEqual, Value: 1},
						{Attribute: "b", Operator: policies.OpEqual, Value: 2},
					},
				},
				attr: policies.MapAttributes{"a": 1, "b": 2},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when any child false should short-circuit and return false",
			input: input{
				pc: policies.PolicyCondition{
					Operator: policies.OpAnd,
					Conditions: []policies.PolicyCondition{
						{Attribute: "a", Operator: policies.OpEqual, Value: 1},
						{Attribute: "b", Operator: policies.OpEqual, Value: 3},
					},
				},
				attr: policies.MapAttributes{"a": 1, "b": 2},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when a child errors should propagate error",
			input: input{
				pc: policies.PolicyCondition{
					Operator:   policies.OpAnd,
					Conditions: []policies.PolicyCondition{{Attribute: "x", Operator: policies.OpEqual, Value: true}},
				},
				attr: policies.MapAttributes{"a": 1},
			},
			output: output{res: false, err: fmt.Errorf("missing required attribute: x")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "missing required attribute")
				}
			},
		},
		{
			name: "when any child true should return true",
			input: input{
				pc: policies.PolicyCondition{
					Operator: policies.OpOr,
					Conditions: []policies.PolicyCondition{
						{Attribute: "a", Operator: policies.OpEqual, Value: 1},
						{Attribute: "b", Operator: policies.OpEqual, Value: 2},
					},
				},
				attr: policies.MapAttributes{"a": 1, "b": 0},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when all children false should return false",
			input: input{
				pc: policies.PolicyCondition{
					Operator: policies.OpOr,
					Conditions: []policies.PolicyCondition{
						{Attribute: "a", Operator: policies.OpEqual, Value: 1},
						{Attribute: "b", Operator: policies.OpEqual, Value: 2},
					},
				},
				attr: policies.MapAttributes{"a": 0, "b": 0},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when children error and none true should return error",
			input: input{
				pc: policies.PolicyCondition{
					Operator: policies.OpOr,
					Conditions: []policies.PolicyCondition{
						{Attribute: "x", Operator: policies.OpEqual, Value: true},
						{Attribute: "b", Operator: policies.OpEqual, Value: 2},
					},
				},
				attr: policies.MapAttributes{"b": 0},
			},
			output: output{res: false, err: fmt.Errorf("missing required attribute: x")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "missing required attribute")
				}
			},
		},
		{
			name: "when a later child true should return true despite earlier errors",
			input: input{
				pc: policies.PolicyCondition{
					Operator: policies.OpOr,
					Conditions: []policies.PolicyCondition{
						{Attribute: "x", Operator: policies.OpEqual, Value: true},
						{Attribute: "b", Operator: policies.OpEqual, Value: 2},
					},
				},
				attr: policies.MapAttributes{"b": 2},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when child true should return false (not)",
			input: input{
				pc:   policies.PolicyCondition{Operator: policies.OpNot, Conditions: []policies.PolicyCondition{{Attribute: "a", Operator: policies.OpEqual, Value: 1}}},
				attr: policies.MapAttributes{"a": 1},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when child errors should return error (not)",
			input: input{
				pc:   policies.PolicyCondition{Operator: policies.OpNot, Conditions: []policies.PolicyCondition{{Attribute: "x", Operator: policies.OpEqual, Value: true}}},
				attr: policies.MapAttributes{"a": 1},
			},
			output: output{res: false, err: fmt.Errorf("missing required attribute: x")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "missing required attribute")
				}
			},
		},
		{
			name: "when not has invalid arity should return error",
			input: input{
				pc:   policies.PolicyCondition{Operator: policies.OpNot, Conditions: []policies.PolicyCondition{{Attribute: "a"}, {Attribute: "b"}}},
				attr: policies.MapAttributes{"a": 1, "b": 2},
			},
			output: output{res: false, err: fmt.Errorf("not operator requires 1 condition")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "not operator requires 1 condition")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use comparison handler to evaluate leaf conditions within tests
			cmp := native.NewComparisonHandler()
			eval := func(pc policies.PolicyCondition, attr policies.Resolver) (bool, error) {
				return cmp.Eval(pc, attr)
			}

			h := native.NewLogicalHandler(eval)

			var actual output
			actual.res, actual.err = h.Eval(tt.input.pc, tt.input.attr)

			tt.assert(t, tt.output, actual)
		})
	}
}
