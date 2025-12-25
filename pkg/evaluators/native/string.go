package native

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

type stringHandler struct{}

func NewStringHandler() policies.OperatorHandler {
	return &stringHandler{}
}

func (hdl *stringHandler) Eval(cond policies.PolicyCondition, ctx policies.AttributeResolver) (bool, error) {
	attr, ok := ctx.Resolve(cond.Attribute)
	if !ok {
		return false, fmt.Errorf("attribute %q not found", cond.Attribute)
	}

	s, ok1 := attr.(string)
	v, ok2 := cond.Value.(string)
	if !ok1 || !ok2 {
		return false, fmt.Errorf("string operator requires string operands")
	}

	switch cond.Operator {
	case policies.OpContains:
		return strings.Contains(s, v), nil

	case policies.OpNotContains:
		return !strings.Contains(s, v), nil

	case policies.OpStartsWith:
		return strings.HasPrefix(s, v), nil

	case policies.OpEndsWith:
		return strings.HasSuffix(s, v), nil

	case policies.OpMatches:
		re, err := regexp.Compile(v)
		if err != nil {
			return false, err
		}
		return re.MatchString(s), nil

	default:
		return false, fmt.Errorf("unsupported string operator: %s", cond.Operator)
	}
}
