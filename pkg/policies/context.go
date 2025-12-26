package policies

import (
	"reflect"
	"strings"
)

type AttributeResolver interface {
	Resolve(attribute string) (any, bool)
}

type MapAttributes map[string]any

func (m MapAttributes) Resolve(attr string) (any, bool) {
	parts := strings.Split(attr, ".")
	var current any = m

	for _, part := range parts {
		switch curr := current.(type) {
		case map[string]any:
			v, ok := curr[part]
			if !ok {
				return nil, false
			}
			current = v

		default:
			val := reflect.ValueOf(curr)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			if val.Kind() == reflect.Struct {
				field := val.FieldByNameFunc(func(name string) bool {
					return strings.EqualFold(name, part)
				})
				if !field.IsValid() {
					return nil, false
				}
				current = field.Interface()
			} else {
				return nil, false
			}
		}
	}

	return current, true
}
