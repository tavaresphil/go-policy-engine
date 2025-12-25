package policies

type AttributeResolver interface {
	Resolve(attribute string) (any, bool)
}

type MapAttributes map[string]any

func (m MapAttributes) Resolve(attr string) (any, bool) {
	v, ok := m[attr]
	return v, ok
}
