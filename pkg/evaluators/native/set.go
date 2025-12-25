package native

import (
	"fmt"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

type setHandler struct{}

func NewSetHandler() policies.OperatorHandler {
	return &setHandler{}
}

func (hdl *setHandler) Eval(cond policies.PolicyCondition, ctx policies.AttributeResolver) (bool, error) {
	attr, ok := ctx.Resolve(cond.Attribute)
	if !ok {
		return false, fmt.Errorf("attribute %q not found", cond.Attribute)
	}

	slice, ok := cond.Value.([]any)
	if !ok {
		return false, fmt.Errorf("set operator requires array value")
	}

	found := false
	for _, v := range slice {
		if v == attr {
			found = true
			break
		}
	}

	switch cond.Operator {
	case policies.OpIn:
		return found, nil
	case policies.OpNotIn:
		return !found, nil
	default:
		return false, fmt.Errorf("unsupported set operator: %s", cond.Operator)
	}
}
