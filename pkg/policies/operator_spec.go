package policies

// OperatorKind groups operators by their execution semantics (comparison, string,
// temporal, logical, etc.). Handlers may be registered by kind to implement
// behavior shared by several operators.
type OperatorKind int

const (
	KindComparison OperatorKind = iota
	KindRange
	KindSet
	KindString
	KindTemporal
	KindArithmetic
	KindLogical
)

// OperatorSpec describes an operator's kind and its arity constraints (min/max
// arguments). MaxArgs == -1 indicates an unbounded number of arguments.
type OperatorSpec struct {
	Kind    OperatorKind
	MinArgs int
	MaxArgs int // -1 = ilimitado
}

// OperatorSpecOf returns the OperatorSpec for a given Operator and whether it exists.
func OperatorSpecOf(op Operator) (OperatorSpec, bool) {
	spec, ok := operatorRegistry[op]
	return spec, ok
}

var operatorRegistry = map[Operator]OperatorSpec{
	// Comparison operators
	OpEqual:          {Kind: KindComparison, MinArgs: 2, MaxArgs: 2},
	OpNotEqual:       {Kind: KindComparison, MinArgs: 2, MaxArgs: 2},
	OpGreater:        {Kind: KindComparison, MinArgs: 2, MaxArgs: 2},
	OpGreaterOrEqual: {Kind: KindComparison, MinArgs: 2, MaxArgs: 2},
	OpLess:           {Kind: KindComparison, MinArgs: 2, MaxArgs: 2},
	OpLessOrEqual:    {Kind: KindComparison, MinArgs: 2, MaxArgs: 2},

	// Range operators
	OpBetween: {Kind: KindRange, MinArgs: 2, MaxArgs: 3}, // valor, min, max (opcional inclusivo)

	// String operators
	OpContains:    {Kind: KindString, MinArgs: 2, MaxArgs: 2},
	OpNotContains: {Kind: KindString, MinArgs: 2, MaxArgs: 2},
	OpStartsWith:  {Kind: KindString, MinArgs: 2, MaxArgs: 2},
	OpEndsWith:    {Kind: KindString, MinArgs: 2, MaxArgs: 2},
	OpMatches:     {Kind: KindString, MinArgs: 2, MaxArgs: 2},

	// Temporal operators
	OpBefore: {Kind: KindTemporal, MinArgs: 2, MaxArgs: 2},
	OpAfter:  {Kind: KindTemporal, MinArgs: 2, MaxArgs: 2},

	// Arithmetic operators
	OpMod: {Kind: KindArithmetic, MinArgs: 2, MaxArgs: 2},

	// Logical operators
	OpAnd: {Kind: KindLogical, MinArgs: 2, MaxArgs: -1}, // Pode ter múltiplos operandos
	OpOr:  {Kind: KindLogical, MinArgs: 2, MaxArgs: -1}, // Pode ter múltiplos operandos
	OpNot: {Kind: KindLogical, MinArgs: 1, MaxArgs: 1},

	// Set operators
	OpIn:         {Kind: KindSet, MinArgs: 2, MaxArgs: 2}, // elemento, conjunto
	OpNotIn:      {Kind: KindSet, MinArgs: 2, MaxArgs: 2}, // elemento, conjunto
	OpSubset:     {Kind: KindSet, MinArgs: 2, MaxArgs: 2}, // subconjunto, conjunto
	OpNotSubset:  {Kind: KindSet, MinArgs: 2, MaxArgs: 2}, // subconjunto, conjunto
	OpIntersects: {Kind: KindSet, MinArgs: 2, MaxArgs: 2}, // conjunto, conjunto
	OpDisjoint:   {Kind: KindSet, MinArgs: 2, MaxArgs: 2}, // conjunto, conjunto
}
