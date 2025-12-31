package native_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tavaresphil/go-policy-engine/pkg/evaluators/native"
	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

func TestStringHandler_Eval(t *testing.T) {
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
				pc:   policies.PolicyCondition{Attribute: "txt", Operator: policies.OpContains, Value: "a"},
				attr: policies.MapAttributes{},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when contains substring should return true",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "txt", Operator: policies.OpContains, Value: "world"},
				attr: policies.MapAttributes{"txt": "hello world"},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when does not contain substring should return false for OpContains",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "txt", Operator: policies.OpContains, Value: "mars"},
				attr: policies.MapAttributes{"txt": "hello world"},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when not contains should invert result for OpNotContains",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "txt", Operator: policies.OpNotContains, Value: "mars"},
				attr: policies.MapAttributes{"txt": "hello world"},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when starts with should return true for OpStartsWith",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "txt", Operator: policies.OpStartsWith, Value: "hello"},
				attr: policies.MapAttributes{"txt": "hello world"},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when ends with should return true for OpEndsWith",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "txt", Operator: policies.OpEndsWith, Value: "world"},
				attr: policies.MapAttributes{"txt": "hello world"},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when regex matches should return true for OpMatches",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "txt", Operator: policies.OpMatches, Value: "^hello"},
				attr: policies.MapAttributes{"txt": "hello world"},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when regex invalid should return error for OpMatches",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "txt", Operator: policies.OpMatches, Value: "("},
				attr: policies.MapAttributes{"txt": "hello world"},
			},
			output: output{res: false, err: fmt.Errorf("invalid regex")},
			assert: func(t *testing.T, expected, actual output) {
				if assert.Error(t, actual.err) {
					assert.Contains(t, actual.err.Error(), "error parsing")
				}
			},
		},
		{
			name: "when attribute convertible to string should work",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "txt", Operator: policies.OpContains, Value: "23"},
				attr: policies.MapAttributes{"txt": 12345},
			},
			output: output{res: true, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when attribute numeric string doesn't contain value should return false",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "txt", Operator: policies.OpContains, Value: "a"},
				attr: policies.MapAttributes{"txt": 1},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
		{
			name: "when value numeric should be converted and not found",
			input: input{
				pc:   policies.PolicyCondition{Attribute: "txt", Operator: policies.OpContains, Value: 1},
				attr: policies.MapAttributes{"txt": "hello"},
			},
			output: output{res: false, err: nil},
			assert: func(t *testing.T, expected, actual output) {
				assert.NoError(t, actual.err)
				assert.Equal(t, expected.res, actual.res)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdl := native.NewStringHandler()

			var actual output
			actual.res, actual.err = hdl.Eval(tt.input.pc, tt.input.attr)

			tt.assert(t, tt.output, actual)
		})
	}
}
