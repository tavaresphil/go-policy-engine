package policies

// Operator is a string-based identifier for a policy operator. Operators
// are declared as constants in this package and are used by PolicyCondition
// to indicate the comparison or operation to be performed.
type Operator string

// Operators used by PolicyCondition. Each operator describes the operation to
// perform during evaluation (comparison, set membership, logical composition,
// etc.).
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

	// Set
	OpIn         Operator = "in"
	OpNotIn      Operator = "nin"
	OpSubset     Operator = "subset"
	OpNotSubset  Operator = "not_subset"
	OpIntersects Operator = "intersects"
	OpDisjoint   Operator = "disjoint"
)
