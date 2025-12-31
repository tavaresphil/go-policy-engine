package native_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tavaresphil/go-policy-engine/pkg/evaluators/native"
	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

func TestComparisonHandler_Eval(t *testing.T) {
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
				pc: policies.PolicyCondition{
					Attribute: "has_dependency",
					Operator:  policies.OpEqual,
					Value:     true,
				},
				attr: policies.MapAttributes{},
			},
			output: output{
				res: false,
				err: fmt.Errorf("missing required attribute: has_dependency"),
			},
			assert: func(t *testing.T, expected, actual output) {
				assert.Equal(t, expected.err.Error(), actual.err.Error())
			},
		},
		{
			name: "when attribute is true should be equal",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "active", Operator: policies.OpEqual, Value: true},
				attr: policies.MapAttributes{"active": true},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when attribute differs should be not equal",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "active", Operator: policies.OpNotEqual, Value: false},
				attr: policies.MapAttributes{"active": true},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when attribute is greater than value should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "count", Operator: policies.OpGreater, Value: 5},
				attr: policies.MapAttributes{"count": 10},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when attribute is less than value should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "count", Operator: policies.OpLess, Value: 5},
				attr: policies.MapAttributes{"count": 3},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when attribute equals value for greater or equal should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "count", Operator: policies.OpGreaterOrEqual, Value: 5},
				attr: policies.MapAttributes{"count": 5},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when attribute equals value for less or equal should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "count", Operator: policies.OpLessOrEqual, Value: 5},
				attr: policies.MapAttributes{"count": 5},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when float attribute is greater than value should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "value", Operator: policies.OpGreater, Value: 2.7},
				attr: policies.MapAttributes{"value": 3.14},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when string attribute is greater than value should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "name", Operator: policies.OpGreater, Value: "a"},
				attr: policies.MapAttributes{"name": "b"},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when time string attribute is greater than value should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "ts", Operator: policies.OpGreater, Value: "2019-12-31"},
				attr: policies.MapAttributes{"ts": "2020-01-02"},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when unix time attribute is greater than value should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "ts", Operator: policies.OpGreater, Value: 1577836800},
				attr: policies.MapAttributes{"ts": 1609459200},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when types differ for greater should return error",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "count", Operator: policies.OpGreater, Value: "5"},
				attr: policies.MapAttributes{"count": 5},
			},
			output: output{res: false, err: fmt.Errorf("cant compare differente type: %T x %T", 5, "5")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "cant compare differente type")
				}
			},
		},
		{
			name: "when unsupported type is compared should return error",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "arr", Operator: policies.OpGreater, Value: []int{1}},
				attr: policies.MapAttributes{"arr": []int{2}},
			},
			output: output{res: false, err: fmt.Errorf("unsupported type for comparison: slice")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "unsupported type for comparison")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdl := native.NewComparisonHandler()

			var actual output
			actual.res, actual.err = hdl.Eval(tt.input.pc, tt.input.attr)

			tt.assert(t, tt.output, actual)
		})
	}
}
