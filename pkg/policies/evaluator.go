package policies

import "fmt"

type EvaluatorRegisty map[OperatorKind]Engine

type evaluator struct {
	reg EvaluatorRegisty
}

func NewEvaluator(reg EvaluatorRegisty) *evaluator {
	return &evaluator{
		reg: reg,
	}
}

func (e *evaluator) Eval(cond PolicyCondition, ctx AttributeResolver) (bool, error) {
	spec, ok := OperatorSpecOf(cond.Operator)
	if !ok {
		return false, fmt.Errorf("unknown operator: %s", cond.Operator)
	}

	handler, ok := e.reg[spec.Kind]
	if !ok {
		return false, fmt.Errorf("no handler for operator kind %v", spec.Kind)
	}

	return handler.Eval(cond, ctx)
}
