package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// AnyToFloat64 converts common numeric representations to float64.
// Supports ints, uints, floats, numeric strings, []byte, fmt.Stringer and pointers (dereferenced).
// Returns an error for nil values or non-numeric inputs.
func AnyToFloat64(v any) (float64, error) {
	if v == nil {
		return 0, fmt.Errorf("nil value")
	}

	// handle pointers
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return 0, fmt.Errorf("nil pointer")
		}
		return AnyToFloat64(rv.Elem().Interface())
	}

	// fmt.Stringer first: convert to string then parse
	if s, ok := v.(fmt.Stringer); ok {
		return strconv.ParseFloat(strings.TrimSpace(s.String()), 64)
	}

	switch t := v.(type) {
	case float64:
		return t, nil
	case float32:
		return float64(t), nil
	case int:
		return float64(t), nil
	case int8:
		return float64(t), nil
	case int16:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case uint:
		return float64(t), nil
	case uint8:
		return float64(t), nil
	case uint16:
		return float64(t), nil
	case uint32:
		return float64(t), nil
	case uint64:
		return float64(t), nil
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return 0, fmt.Errorf("empty string")
		}
		return strconv.ParseFloat(s, 64)
	case []byte:
		return strconv.ParseFloat(strings.TrimSpace(string(t)), 64)
	default:
		// last resort: try fmt.Sprint and parse
		if s := fmt.Sprint(v); s != "" {
			if f, err := strconv.ParseFloat(strings.TrimSpace(s), 64); err == nil {
				return f, nil
			}
		}
		return 0, fmt.Errorf("not numeric: %T", v)
	}
}
