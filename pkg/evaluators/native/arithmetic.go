package native

import (
	"fmt"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

type arithmeticHandler struct{}

func NewArithmeticHandler() policies.OperatorHandler {
	return &arithmeticHandler{}
}

func (hdl *arithmeticHandler) Eval(cond policies.PolicyCondition, ctx policies.AttributeResolver) (bool, error) {
	attr, ok := ctx.Resolve(cond.Attribute)
	if !ok {
		return false, fmt.Errorf("attribute %q not found", cond.Attribute)
	}

	switch cond.Operator {
	case policies.OpMod:
		a, ok1 := toInt(attr)
		b, ok2 := toInt(cond.Value)
		if !ok1 || !ok2 {
			return false, fmt.Errorf("mod operator requires integer operands")
		}
		if b == 0 {
			return false, fmt.Errorf("mod by zero")
		}
		return a%b == 0, nil

	default:
		return false, fmt.Errorf("unsupported arithmetic operator: %s", cond.Operator)
	}
}

func toInt(v any) (int, bool) {
	switch t := v.(type) {
	case int:
		return t, true
	case int64:
		return int(t), true
	case float64:
		return int(t), true
	default:
		return 0, false
	}
}
