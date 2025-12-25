package policies

type Operator string

const (
	// Comparison
	OpEqual          Operator = "eq"
	OpNotEqual       Operator = "neq"
	OpGreater        Operator = "gt"
	OpGreaterOrEqual Operator = "gte"
	OpLess           Operator = "lt"
	OpLessOrEqual    Operator = "lte"

	// Range
	OpBetween Operator = "between"

	// Set
	OpIn    Operator = "in"
	OpNotIn Operator = "nin"

	// String
	OpContains    Operator = "contains"
	OpNotContains Operator = "not_contains"
	OpStartsWith  Operator = "starts_with"
	OpEndsWith    Operator = "ends_with"
	OpMatches     Operator = "matches"

	// Temporal
	OpBefore Operator = "before"
	OpAfter  Operator = "after"

	// Arithmetic
	OpMod Operator = "mod"

	// Logical
	OpAnd Operator = "and"
	OpOr  Operator = "or"
	OpNot Operator = "not"
)
