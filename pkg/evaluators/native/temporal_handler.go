package native

import (
	"fmt"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
	"github.com/tavaresphil/go-policy-engine/pkg/utils"
)

type TemporalHandler struct{}

func NewTemporalHandler() OperatorHandler {
	return &TemporalHandler{}
}
func (h *TemporalHandler) Eval(pc policies.PolicyCondition, attr policies.Resolver) (bool, error) {
	attrVal, ok := attr.Resolve(pc.Attribute)
	if !ok {
		return false, fmt.Errorf("missing required attribute: %s", pc.Attribute)
	}

	attrTime, err := utils.AnyToTime(attrVal)
	if err != nil {
		return false, err
	}

	valueTime, err := utils.AnyToTime(pc.Value)
	if err != nil {
		return false, err
	}
	switch pc.Operator {
	case policies.OpBefore:
		return attrTime.Before(valueTime), nil
	case policies.OpAfter:
		return attrTime.After(valueTime), nil
	default:
		return false, fmt.Errorf("unsupported temporal operator: %s", pc.Operator)
	}
}
