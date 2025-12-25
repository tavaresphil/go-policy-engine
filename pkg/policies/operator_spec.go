package policies

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

type OperatorSpec struct {
	Kind    OperatorKind
	MinArgs int
	MaxArgs int // -1 = ilimitado
}

var operatorRegistry = map[Operator]OperatorSpec{
	// Comparison
	OpEqual:          {Kind: KindComparison, MinArgs: 1, MaxArgs: 1},
	OpNotEqual:       {Kind: KindComparison, MinArgs: 1, MaxArgs: 1},
	OpGreater:        {Kind: KindComparison, MinArgs: 1, MaxArgs: 1},
	OpGreaterOrEqual: {Kind: KindComparison, MinArgs: 1, MaxArgs: 1},
	OpLess:           {Kind: KindComparison, MinArgs: 1, MaxArgs: 1},
	OpLessOrEqual:    {Kind: KindComparison, MinArgs: 1, MaxArgs: 1},

	// Range
	OpBetween: {Kind: KindRange, MinArgs: 2, MaxArgs: 2},

	// Set
	OpIn:    {Kind: KindSet, MinArgs: 1, MaxArgs: 1},
	OpNotIn: {Kind: KindSet, MinArgs: 1, MaxArgs: 1},

	// String
	OpContains:    {Kind: KindString, MinArgs: 1, MaxArgs: 1},
	OpNotContains: {Kind: KindString, MinArgs: 1, MaxArgs: 1},
	OpStartsWith:  {Kind: KindString, MinArgs: 1, MaxArgs: 1},
	OpEndsWith:    {Kind: KindString, MinArgs: 1, MaxArgs: 1},
	OpMatches:     {Kind: KindString, MinArgs: 1, MaxArgs: 1},

	// Temporal
	OpBefore: {Kind: KindTemporal, MinArgs: 1, MaxArgs: 1},
	OpAfter:  {Kind: KindTemporal, MinArgs: 1, MaxArgs: 1},

	// Arithmetic
	OpMod: {Kind: KindArithmetic, MinArgs: 1, MaxArgs: 1},

	// Logical
	OpAnd: {Kind: KindLogical, MinArgs: 1, MaxArgs: -1},
	OpOr:  {Kind: KindLogical, MinArgs: 1, MaxArgs: -1},
	OpNot: {Kind: KindLogical, MinArgs: 1, MaxArgs: 1},
}

func OperatorSpecOf(op Operator) (OperatorSpec, bool) {
	spec, ok := operatorRegistry[op]
	return spec, ok
}
