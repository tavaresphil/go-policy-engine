package utils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tavaresphil/go-policy-engine/pkg/utils"
)

type myStringer struct{ v string }

func (m myStringer) String() string { return m.v }

func TestAnyToString(t *testing.T) {
	tm := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name string
		in   any
		out  string
		err  bool
	}{
		{"string", "hello", "hello", false},
		{"bytes", []byte("abc"), "abc", false},
		{"stringer", myStringer{"x"}, "x", false},
		{"bool", true, "true", false},
		{"int", 123, "123", false},
		{"float", 3.14, "3.14", false},
		{"time", tm, "2020-01-02T03:04:05Z", false},
		{"time ptr", &tm, "2020-01-02T03:04:05Z", false},
		{"nil ptr", (*int)(nil), "", true},
		{"nil", nil, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := utils.AnyToString(tt.in)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.out, out)
			}
		})
	}
}
