package policies

// OperatorHandler evaluates a PolicyCondition against an attribute context.
// Implementations provide handlers for operator kinds or specific operators.
type OperatorHandler interface {
	Eval(cond PolicyCondition, ctx Resolver) (bool, error)
}
