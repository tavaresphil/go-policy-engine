package policies

import (
	"reflect"
	"strings"
)

// Resolver resolves attribute values by name. The attribute name may be a dotted path
// (for example "user.name") to traverse nested maps or structs. The boolean return
// value indicates whether the attribute was present.
//
// Implementations should avoid panicking and return (nil, false) when the attribute
// cannot be resolved.
type Resolver interface {
	Resolve(attribute string) (any, bool)
}

// MapAttributes is a map-based implementation of Resolver that supports dotted paths
// to traverse nested maps and structs.
type MapAttributes map[string]any

// Resolve attempts to find the attribute identified by attr. Supports dotted
// paths to traverse nested maps and exported struct fields (case-insensitive).
// Returns (value, true) if found or (nil, false) otherwise.
func (m MapAttributes) Resolve(attr string) (any, bool) {
	parts := strings.Split(attr, ".")
	var current any = m

	for _, part := range parts {
		switch curr := current.(type) {
		case MapAttributes:
			v, ok := curr[part]
			if !ok {
				return nil, false
			}
			current = v
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
