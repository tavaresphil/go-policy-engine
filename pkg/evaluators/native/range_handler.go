package native

import (
	"fmt"
	"reflect"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
	"github.com/tavaresphil/go-policy-engine/pkg/utils"
)

type RangeHandler struct{}

func NewRangeHandler() OperatorHandler {
	return &RangeHandler{}
}

// parseBetweenValue accepts the PolicyCondition.Value which can be:
// - a slice/array with 2 elements: [min, max]
// - a slice/array with 3 elements: [min, max, inclusive(bool)]
// - a map[string]any with keys "min","max","inclusive"
func parseBetweenValue(v any) (min any, max any, inclusive bool, err error) {
	inclusive = true // default inclusive
	if v == nil {
		return nil, nil, false, fmt.Errorf("between requires min and max")
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		if rv.Len() < 2 || rv.Len() > 3 {
			return nil, nil, false, fmt.Errorf("between requires 2 or 3 args: min, max, (inclusive)")
		}
		min = rv.Index(0).Interface()
		max = rv.Index(1).Interface()
		if rv.Len() == 3 {
			incVal := rv.Index(2).Interface()
			b, ok := incVal.(bool)
			if !ok {
				return nil, nil, false, fmt.Errorf("inclusive flag must be a boolean")
			}
			inclusive = b
		}
		return
	case reflect.Map:
		// try map[string]any
		m, ok := v.(map[string]any)
		if !ok {
			return nil, nil, false, fmt.Errorf("unsupported between value map type: %T", v)
		}
		var okMin, okMax bool
		min, okMin = m["min"]
		max, okMax = m["max"]
		if !okMin || !okMax {
			return nil, nil, false, fmt.Errorf("between requires min and max in map form")
		}
		if inc, ok := m["inclusive"]; ok {
			b, ok := inc.(bool)
			if !ok {
				return nil, nil, false, fmt.Errorf("inclusive flag must be a boolean")
			}
			inclusive = b
		}
		return
	default:
		return nil, nil, false, fmt.Errorf("unsupported between value type: %v", rv.Kind())
	}
}

func (h *RangeHandler) Eval(pc policies.PolicyCondition, attr policies.Resolver) (bool, error) {
	if pc.Operator != policies.OpBetween {
		return false, fmt.Errorf("unsupported range operator: %s", pc.Operator)
	}

	attrVal, ok := attr.Resolve(pc.Attribute)
	if !ok {
		return false, fmt.Errorf("missing required attribute: %s", pc.Attribute)
	}

	min, max, inclusive, err := parseBetweenValue(pc.Value)
	if err != nil {
		return false, err
	}

	// Try parsing everything as time
	if at, bt, errt := ParseBothAsTime(attrVal, min); errt == nil {
		if mt, errt2 := utils.AnyToTime(max); errt2 == nil {
			if inclusive {
				return (at.Equal(bt) || at.After(bt)) && (at.Equal(mt) || at.Before(mt)), nil
			}
			return at.After(bt) && at.Before(mt), nil
		}
		// fallthrough if max isn't time
	}

	// Try numeric coercion using utils.AnyToFloat64 which handles numbers and numeric strings
	if af, errA := utils.AnyToFloat64(attrVal); errA == nil {
		if mf, errB := utils.AnyToFloat64(min); errB == nil {
			if Mf, errC := utils.AnyToFloat64(max); errC == nil {
				// if min > max, return false (not in range)
				if mf > Mf {
					return false, nil
				}
				if inclusive {
					return af >= mf && af <= Mf, nil
				}
				return af > mf && af < Mf, nil
			}
		}
	}

	// strings using AnyToString fallback for convertible types (includes fmt.Stringer, numbers, bytes, etc.)
	if sa, errA := utils.AnyToString(attrVal); errA == nil {
		if sm, errB := utils.AnyToString(min); errB == nil {
			if sM, errC := utils.AnyToString(max); errC == nil {
				// lexicographic check
				if sm > sM {
					return false, nil
				}
				if inclusive {
					return sa >= sm && sa <= sM, nil
				}
				return sa > sm && sa < sM, nil
			}
		}
	}

	// final failure
	return false, fmt.Errorf("mismatched types for between: %T, %T, %T", attrVal, min, max)
}
