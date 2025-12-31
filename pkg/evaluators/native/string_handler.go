package native

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
	"github.com/tavaresphil/go-policy-engine/pkg/utils"
)

type stringHandler struct{}

func NewStringHandler() OperatorHandler {
	return &stringHandler{}
}

func (h *stringHandler) Eval(pc policies.PolicyCondition, attr policies.Resolver) (bool, error) {
	attrVal, ok := attr.Resolve(pc.Attribute)
	if !ok {
		return false, nil
	}

	// convert attribute and value to string using helper
	as, err := utils.AnyToString(attrVal)
	if err != nil {
		return false, err
	}

	vs, err := utils.AnyToString(pc.Value)
	if err != nil {
		return false, err
	}

	switch pc.Operator {
	case policies.OpContains:
		return strings.Contains(as, vs), nil
	case policies.OpNotContains:
		return !strings.Contains(as, vs), nil
	case policies.OpStartsWith:
		return strings.HasPrefix(as, vs), nil
	case policies.OpEndsWith:
		return strings.HasSuffix(as, vs), nil
	case policies.OpMatches:
		re, err := regexp.Compile(vs)
		if err != nil {
			return false, err
		}
		return re.MatchString(as), nil
	default:
		return false, fmt.Errorf("unsupported string operator: %s", pc.Operator)
	}
}
