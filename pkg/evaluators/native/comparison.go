package native

import (
	"fmt"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

type comparisonHandler struct{}

func NewComparisonHandler() policies.OperatorHandler {
	return &comparisonHandler{}
}

func (hdl *comparisonHandler) Eval(cond policies.PolicyCondition, ctx policies.AttributeResolver) (bool, error) {
	left, ok := ctx.Resolve(cond.Attribute)
	if !ok {
		return false, fmt.Errorf("attribute %q not found", cond.Attribute)
	}

	switch l := left.(type) {
	case int, int64, float64:
		return compareNumbers(cond.Operator, l, cond.Value)

	case string:
		r, ok := cond.Value.(string)
		if !ok {
			return false, fmt.Errorf("string comparison requires string value")
		}
		return compareStrings(cond.Operator, l, r)

	default:
		return false, fmt.Errorf("unsupported comparison type: %T", left)
	}
}

func compareNumbers(op policies.Operator, a any, b any) (bool, error) {
	x, ok1 := toFloat(a)
	y, ok2 := toFloat(b)
	if !ok1 || !ok2 {
		return false, fmt.Errorf("numeric comparison requires numbers")
	}

	switch op {
	case policies.OpEqual:
		return x == y, nil
	case policies.OpNotEqual:
		return x != y, nil
	case policies.OpGreater:
		return x > y, nil
	case policies.OpGreaterOrEqual:
		return x >= y, nil
	case policies.OpLess:
		return x < y, nil
	case policies.OpLessOrEqual:
		return x <= y, nil
	default:
		return false, fmt.Errorf("unsupported comparison operator: %s", op)
	}
}

func compareStrings(op policies.Operator, a, b string) (bool, error) {
	switch op {
	case policies.OpEqual:
		return a == b, nil
	case policies.OpNotEqual:
		return a != b, nil
	case policies.OpGreater:
		return a > b, nil
	case policies.OpGreaterOrEqual:
		return a >= b, nil
	case policies.OpLess:
		return a < b, nil
	case policies.OpLessOrEqual:
		return a <= b, nil
	default:
		return false, fmt.Errorf("unsupported string comparison operator: %s", op)
	}
}

func toFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case float64:
		return t, true
	default:
		return 0, false
	}
}
