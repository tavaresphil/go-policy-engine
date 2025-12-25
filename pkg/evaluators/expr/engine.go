package expr

import (
	"fmt"

	"github.com/expr-lang/expr"
	exprlang "github.com/expr-lang/expr"
	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

type engine struct {
	builder ExprBuilder
}

func NewEngine() policies.Engine {
	return &engine{
		builder: NewExprBuilder(),
	}
}

func (e *engine) Eval(cond policies.PolicyCondition, ctx policies.AttributeResolver) (bool, error) {
	if err := cond.Validate(); err != nil {
		return false, fmt.Errorf("invalid condition: %w", err)
	}

	exprStr, err := e.builder.Build(cond)
	if err != nil {
		return false, fmt.Errorf("failed to build expression: %w", err)
	}

	program, err := exprlang.Compile(exprStr, exprlang.Env(ctx), expr.AsBool())
	if err != nil {
		return false, fmt.Errorf("failed compile expression: %w", err)
	}

	res, err := expr.Run(program, ctx)
	if err != nil {
		return false, fmt.Errorf("failed run expression: %w", err)
	}

	matches, _ := res.(bool)
	return matches, nil
}
