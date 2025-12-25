package policies

// Engine evaluates policy conditions.
type Engine interface {
	Eval(cond PolicyCondition, ctx AttributeResolver) (bool, error)
}
