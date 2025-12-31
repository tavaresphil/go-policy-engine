package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// AnyToString attempts to convert common Go types to their string representation.
// It handles strings, []byte, fmt.Stringer, booleans, numbers, and time.Time.
// For pointers it will dereference and attempt conversion; nil pointers return an error.
func AnyToString(v any) (string, error) {
	if v == nil {
		return "", fmt.Errorf("nil value")
	}

	// handle time types before the generic Stringer to get RFC3339 format
	switch t := v.(type) {
	case time.Time:
		return t.Format(time.RFC3339), nil
	case *time.Time:
		if t == nil {
			return "", fmt.Errorf("nil time pointer")
		}
		return t.Format(time.RFC3339), nil
	}

	// handle fmt.Stringer implicitly for other types
	if s, ok := v.(fmt.Stringer); ok {
		return s.String(), nil
	}

	switch t := v.(type) {
	case string:
		return t, nil
	case []byte:
		return string(t), nil
	case bool:
		return strconv.FormatBool(t), nil
	case int, int8, int16, int32, int64:
		return fmt.Sprint(t), nil
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprint(t), nil
	case float32, float64:
		return fmt.Sprint(t), nil
	}

	// handle pointers generically
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "", fmt.Errorf("nil pointer")
		}
		return AnyToString(rv.Elem().Interface())
	}

	// Fallback to sprint for other types (maps, slices, structs)
	return fmt.Sprint(v), nil
}
