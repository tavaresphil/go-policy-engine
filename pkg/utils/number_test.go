package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tavaresphil/go-policy-engine/pkg/utils"
)

func TestAnyToFloat64(t *testing.T) {
	tests := []struct {
		name string
		in   any
		out  float64
		err  bool
	}{
		{"int", 123, 123, false},
		{"int64", int64(123), 123, false},
		{"uint", uint(5), 5, false},
		{"float32", float32(3.14), 3.14, false}, // float32 precision checked with delta
		{"float64", 3.14, 3.14, false},
		{"string int", "42", 42, false},
		{"string float", "  2.5 ", 2.5, false},
		{"bytes", []byte("7.5"), 7.5, false},
		{"numeric stringer", struct{ s string }{"8"}, 8, true},
		{"non-numeric string", "abc", 0, true},
		{"nil", nil, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := utils.AnyToFloat64(tt.in)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.name == "float32" {
					assert.InDelta(t, tt.out, out, 1e-6)
				} else {
					assert.Equal(t, tt.out, out)
				}
			}
		})
	}
}
