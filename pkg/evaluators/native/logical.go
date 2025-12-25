package native

import (
	"fmt"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

type logicalHandler struct {
	eval EvalFunc
}

func NewLogicalHandler(eval EvalFunc) policies.OperatorHandler {
	return &logicalHandler{eval: eval}
}

func (hdl *logicalHandler) Eval(cond policies.PolicyCondition, ctx policies.AttributeResolver) (bool, error) {

	switch cond.Operator {
	case policies.OpAnd:
		for _, c := range cond.Conditions {
			ok, err := hdl.eval(c, ctx)
			if err != nil || !ok {
				return false, err
			}
		}
		return true, nil

	case policies.OpOr:
		for _, c := range cond.Conditions {
			ok, err := hdl.eval(c, ctx)
			if err == nil && ok {
				return true, nil
			}
		}
		return false, nil

	case policies.OpNot:
		if len(cond.Conditions) != 1 {
			return false, fmt.Errorf("not expects exactly one condition")
		}
		ok, err := hdl.eval(cond.Conditions[0], ctx)
		if err != nil {
			return false, err
		}
		return !ok, nil

	default:
		return false, fmt.Errorf("unsupported logical operator %s", cond.Operator)
	}
}
