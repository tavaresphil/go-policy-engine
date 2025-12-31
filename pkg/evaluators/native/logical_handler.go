package native

import (
	"fmt"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

type Eval func(pc policies.PolicyCondition, attr policies.Resolver) (bool, error)

type LogicalHandler struct {
	eval Eval
}

func NewLogicalHandler(eval Eval) OperatorHandler {
	return &LogicalHandler{eval: eval}
}

func (h *LogicalHandler) Eval(pc policies.PolicyCondition, attr policies.Resolver) (bool, error) {
	switch pc.Operator {
	case policies.OpAnd:
		for _, cond := range pc.Conditions {
			ok, err := h.eval(cond, attr)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
		}
		return true, nil
	case policies.OpOr:
		var lastErr error
		for _, cond := range pc.Conditions {
			ok, err := h.eval(cond, attr)
			if err != nil {
				lastErr = err
				continue
			}
			if ok {
				return true, nil
			}
		}
		if lastErr != nil {
			return false, lastErr
		}
		return false, nil
	case policies.OpNot:
		if len(pc.Conditions) != 1 {
			return false, fmt.Errorf("not operator requires 1 condition")
		}
		ok, err := h.eval(pc.Conditions[0], attr)
		if err != nil {
			return false, err
		}
		return !ok, nil
	default:
		return false, fmt.Errorf("unsupported logical operator: %s", pc.Operator)
	}
}
