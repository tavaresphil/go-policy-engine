package native

import (
	"fmt"
	"reflect"
	"time"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
	"github.com/tavaresphil/go-policy-engine/pkg/utils"
)

type ComparisonHandler struct{}

func NewComparisonHandler() OperatorHandler {
	return &ComparisonHandler{}
}

func (h *ComparisonHandler) Eval(pc policies.PolicyCondition, attr policies.Resolver) (bool, error) {
	attrVal, ok := attr.Resolve(pc.Attribute)
	if !ok {
		return false, fmt.Errorf("missing required attribute: %s", pc.Attribute)
	}

	switch pc.Operator {
	case policies.OpEqual:
		return equal(attrVal, pc.Value), nil
	case policies.OpNotEqual:
		return !equal(attrVal, pc.Value), nil
	case policies.OpGreater:
		return greater(attrVal, pc.Value)
	case policies.OpGreaterOrEqual:
		gt, err := greater(attrVal, pc.Value)
		if err != nil {
			return false, err
		}
		return gt || equal(attrVal, pc.Value), nil
	case policies.OpLess:
		gt, err := greater(attrVal, pc.Value)
		if err != nil {
			return false, err
		}
		eq := equal(attrVal, pc.Value)
		return !gt && !eq, nil
	case policies.OpLessOrEqual:
		gt, err := greater(attrVal, pc.Value)
		if err != nil {
			return false, err
		}
		return !gt || equal(attrVal, pc.Value), nil
	default:
		return false, fmt.Errorf("unsupported comparison operator: %s", pc.Operator)
	}
}

func equal(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

func greater(a, b any) (bool, error) {
	if atime, btime, err := ParseBothAsTime(a, b); err == nil {
		return atime.After(btime), nil
	}

	aval := reflect.ValueOf(a)
	bval := reflect.ValueOf(b)

	if aval.Type() != bval.Type() {
		return false, fmt.Errorf("cant compare differente type: %T x %T", a, b)
	}

	switch aval.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return aval.Int() > bval.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return aval.Uint() > bval.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return aval.Float() > bval.Float(), nil
	case reflect.String:
		return aval.String() > bval.String(), nil
	default:
		return false, fmt.Errorf("unsupported type for comparison: %v", aval.Kind())
	}
}

func ParseBothAsTime(a, b any) (time.Time, time.Time, error) {
	var atime, btime time.Time
	var err error

	atime, err = utils.AnyToTime(a)
	if err != nil {
		return atime, btime, err
	}

	btime, err = utils.AnyToTime(b)
	if err != nil {
		return atime, btime, err
	}

	return atime, btime, nil
}
