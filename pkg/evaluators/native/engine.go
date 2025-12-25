package native

import (
	"fmt"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

type EvalFunc func(cond policies.PolicyCondition, ctx policies.AttributeResolver) (bool, error)

type Engine struct {
	handlers map[policies.Operator]policies.OperatorHandler
}

var _ policies.Engine = (*Engine)(nil)

func NewEngine() *Engine {
	e := &Engine{
		handlers: make(map[policies.Operator]policies.OperatorHandler),
	}

	e.registerHandlers()
	return e
}

func (e *Engine) Eval(cond policies.PolicyCondition, ctx policies.AttributeResolver) (bool, error) {
	handler, ok := e.handlers[cond.Operator]
	if !ok {
		return false, fmt.Errorf("no handler registered for operator %s", cond.Operator)
	}

	return handler.Eval(cond, ctx)
}
