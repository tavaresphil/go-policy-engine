package native

import (
	"fmt"
	"time"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

type temporalHandler struct{}

func NewTemporalHandler() policies.OperatorHandler {
	return &temporalHandler{}
}

func (hdl *temporalHandler) Eval(cond policies.PolicyCondition, ctx policies.AttributeResolver) (bool, error) {
	attr, ok := ctx.Resolve(cond.Attribute)
	if !ok {
		return false, fmt.Errorf("attribute %q not found", cond.Attribute)
	}

	left, ok1 := attr.(time.Time)
	right, ok2 := cond.Value.(time.Time)
	if !ok1 || !ok2 {
		return false, fmt.Errorf("temporal operator requires time.Time values")
	}

	switch cond.Operator {
	case policies.OpBefore:
		return left.Before(right), nil
	case policies.OpAfter:
		return left.After(right), nil
	default:
		return false, fmt.Errorf("unsupported temporal operator: %s", cond.Operator)
	}
}
