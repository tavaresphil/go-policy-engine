package native

import (
	"fmt"
	"reflect"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

type SetHandler struct{}

func NewSetHandler() OperatorHandler {
	return &SetHandler{}
}

func (h *SetHandler) Eval(pc policies.PolicyCondition, attr policies.Resolver) (bool, error) {
	attrVal, ok := attr.Resolve(pc.Attribute)
	if !ok {
		return false, nil
	}
	switch pc.Operator {
	case policies.OpIn:
		return contains(pc.Value, attrVal)
	case policies.OpNotIn:
		ok, err := contains(pc.Value, attrVal)
		if err != nil {
			return false, err
		}
		return !ok, nil
	case policies.OpSubset:
		return isSubset(attrVal, pc.Value)
	case policies.OpNotSubset:
		subset, err := isSubset(attrVal, pc.Value)
		return !subset, err
	case policies.OpIntersects:
		return intersects(attrVal, pc.Value)
	case policies.OpDisjoint:
		inter, err := intersects(attrVal, pc.Value)
		return !inter, err
	default:
		return false, fmt.Errorf("unsupported set operator: %s", pc.Operator)
	}
}

func contains(set any, item any) (bool, error) {
	setVal := reflect.ValueOf(set)
	switch setVal.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < setVal.Len(); i++ {
			if reflect.DeepEqual(setVal.Index(i).Interface(), item) {
				return true, nil
			}
		}
		return false, nil
	case reflect.Map:
		mapKeys := setVal.MapKeys()
		for _, key := range mapKeys {
			if reflect.DeepEqual(key.Interface(), item) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, nil
	}
}

func isSubset(subset any, set any) (bool, error) {
	subsetVal := reflect.ValueOf(subset)
	setVal := reflect.ValueOf(set)

	if subsetVal.Kind() != reflect.Slice && subsetVal.Kind() != reflect.Array {
		return false, fmt.Errorf("subset must be a slice or array, got %v", subsetVal.Kind())
	}
	if setVal.Kind() != reflect.Slice && setVal.Kind() != reflect.Array {
		return false, fmt.Errorf("superset must be a slice or array, got %v", setVal.Kind())
	}

	setMap := make(map[any]bool)
	for i := 0; i < setVal.Len(); i++ {
		setMap[setVal.Index(i).Interface()] = true
	}

	for i := 0; i < subsetVal.Len(); i++ {
		if !setMap[subsetVal.Index(i).Interface()] {
			return false, nil
		}
	}

	return true, nil
}

func intersects(setA, setB any) (bool, error) {
	setAVal := reflect.ValueOf(setA)
	setBVal := reflect.ValueOf(setB)

	if setAVal.Kind() != reflect.Slice && setAVal.Kind() != reflect.Array {
		return false, fmt.Errorf("first set must be a slice or array")
	}
	if setBVal.Kind() != reflect.Slice && setBVal.Kind() != reflect.Array {
		return false, fmt.Errorf("second set must be a slice or array")
	}

	var smaller, larger reflect.Value
	if setAVal.Len() < setBVal.Len() {
		smaller = setAVal
		larger = setBVal
	} else {
		smaller = setBVal
		larger = setAVal
	}

	smallerMap := make(map[any]bool)
	for i := 0; i < smaller.Len(); i++ {
		smallerMap[smaller.Index(i).Interface()] = true
	}

	for i := 0; i < larger.Len(); i++ {
		if smallerMap[larger.Index(i).Interface()] {
			return true, nil
		}
	}

	return false, nil
}
