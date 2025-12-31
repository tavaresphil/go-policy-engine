package native

import (
	"fmt"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

type OperatorHandler interface {
	Eval(pc policies.PolicyCondition, attr policies.Resolver) (bool, error)
}

type NativeEngine struct {
	handlers map[policies.OperatorKind]OperatorHandler
}

func NewNativeEngine() policies.Engine {
	handlers := make(map[policies.OperatorKind]OperatorHandler)
	handlers[policies.KindArithmetic] = NewArithmeticHandler()
	handlers[policies.KindComparison] = NewComparisonHandler()
	handlers[policies.KindString] = NewStringHandler()
	handlers[policies.KindSet] = NewSetHandler()
	handlers[policies.KindRange] = NewRangeHandler()
	handlers[policies.KindTemporal] = NewTemporalHandler()

	eng := &NativeEngine{handlers: handlers}
	eng.handlers[policies.KindLogical] = NewLogicalHandler(eng.Eval)
	return eng
}

func (e *NativeEngine) Eval(pc policies.PolicyCondition, attr policies.Resolver) (bool, error) {
	spec, ok := policies.OperatorSpecOf(pc.Operator)
	if !ok {
		return false, fmt.Errorf("unknown operator: %s", pc.Operator)
	}

	handler, ok := e.handlers[spec.Kind]
	if !ok {
		return false, fmt.Errorf("unsupported operator kind: %v", spec.Kind)
	}

	return handler.Eval(pc, attr)
}
