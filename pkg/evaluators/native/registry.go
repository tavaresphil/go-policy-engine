package native

import "github.com/tavaresphil/go-policy-engine/pkg/policies"

func (e *Engine) registerHandlers() {
	// logical
	logical := NewLogicalHandler(e.Eval)
	e.bind(policies.OpAnd, logical)
	e.bind(policies.OpOr, logical)
	e.bind(policies.OpNot, logical)

	// comparison
	comp := NewComparisonHandler()
	e.bind(
		policies.OpEqual,
		policies.OpNotEqual,
		policies.OpGreater,
		policies.OpGreaterOrEqual,
		policies.OpLess,
		policies.OpLessOrEqual,
		comp,
	)

	// string
	str := NewStringHandler()
	e.bind(
		policies.OpContains,
		policies.OpNotContains,
		policies.OpStartsWith,
		policies.OpEndsWith,
		policies.OpMatches,
		str,
	)

	// set
	set := NewSetHandler()
	e.bind(policies.OpIn, policies.OpNotIn, set)

	// temporal
	tmp := NewTemporalHandler()
	e.bind(policies.OpBefore, policies.OpAfter, tmp)

	// arithmetic
	ar := NewArithmeticHandler()
	e.bind(policies.OpMod, ar)
}

func (e *Engine) bind(ops ...any) {
	handler := ops[len(ops)-1].(policies.OperatorHandler)
	for _, op := range ops[:len(ops)-1] {
		e.handlers[op.(policies.Operator)] = handler
	}
}
