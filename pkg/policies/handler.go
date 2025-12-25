package policies

type OperatorHandler interface {
	Eval(cond PolicyCondition, ctx AttributeResolver) (bool, error)
}
