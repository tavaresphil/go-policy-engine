package native_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tavaresphil/go-policy-engine/pkg/evaluators/native"
	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

func TestSetHandler_Eval(t *testing.T) {
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
			name: "when attribute not found should return false",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "x", Operator: policies.OpIn, Value: []int{1, 2}},
				attr: policies.MapAttributes{},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when value is in slice should return true for OpIn",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "item", Operator: policies.OpIn, Value: []int{1, 2}},
				attr: policies.MapAttributes{"item": 1},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when value is not in slice should return false for OpIn",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "item", Operator: policies.OpIn, Value: []int{2, 3}},
				attr: policies.MapAttributes{"item": 1},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when value is map key should return true for OpIn",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "item", Operator: policies.OpIn, Value: map[int]bool{1: true}},
				attr: policies.MapAttributes{"item": 1},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when value not in map keys should return false for OpIn",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "item", Operator: policies.OpIn, Value: map[int]bool{2: true}},
				attr: policies.MapAttributes{"item": 1},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when not in should invert result for OpNotIn",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "item", Operator: policies.OpNotIn, Value: []int{2, 3}},
				attr: policies.MapAttributes{"item": 1},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when attribute is subset should return true for OpSubset",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "subset", Operator: policies.OpSubset, Value: []int{1, 2, 3}},
				attr: policies.MapAttributes{"subset": []int{1, 2}},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when attribute is not subset should return false for OpSubset",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "subset", Operator: policies.OpSubset, Value: []int{2, 3}},
				attr: policies.MapAttributes{"subset": []int{1, 2}},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when subset arg not slice should return error for OpSubset",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "subset", Operator: policies.OpSubset, Value: []int{1, 2}},
				attr: policies.MapAttributes{"subset": 1},
			},
			output: output{res: false, err: fmt.Errorf("subset must be a slice or array")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "subset must be a slice or array")
				}
			},
		},
		{
			name: "when superset arg not slice should return error for OpSubset",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "subset", Operator: policies.OpSubset, Value: 1},
				attr: policies.MapAttributes{"subset": []int{1}},
			},
			output: output{res: false, err: fmt.Errorf("superset must be a slice or array")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "superset must be a slice or array")
				}
			},
		},
		{
			name: "when sets intersect should return true for OpIntersects",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "a", Operator: policies.OpIntersects, Value: []int{3, 4}},
				attr: policies.MapAttributes{"a": []int{1, 3}},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when sets do not intersect should return false for OpIntersects",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "a", Operator: policies.OpIntersects, Value: []int{4, 5}},
				attr: policies.MapAttributes{"a": []int{1, 3}},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when first arg not slice should return error for OpIntersects",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "a", Operator: policies.OpIntersects, Value: []int{1}},
				attr: policies.MapAttributes{"a": 1},
			},
			output: output{res: false, err: fmt.Errorf("first set must be a slice or array")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "first set must be a slice or array")
				}
			},
		},
		{
			name: "when second arg not slice should return error for OpIntersects",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "a", Operator: policies.OpIntersects, Value: 1},
				attr: policies.MapAttributes{"a": []int{1}},
			},
			output: output{res: false, err: fmt.Errorf("second set must be a slice or array")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "second set must be a slice or array")
				}
			},
		},
		{
			name: "when sets are disjoint should return true for OpDisjoint",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "a", Operator: policies.OpDisjoint, Value: []int{4, 5}},
				attr: policies.MapAttributes{"a": []int{1, 3}},
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
			hdl := native.NewSetHandler()

			var actual output
			actual.res, actual.err = hdl.Eval(tt.input.pc, tt.input.attr)

			tt.assert(t, tt.output, actual)
		})
	}
}
