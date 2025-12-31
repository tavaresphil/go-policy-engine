package expr

import (
	"fmt"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

// ExprBuilder translates a PolicyCondition into an expr-lang expression fragment.
// Implementations produce a string expression fragment that can be combined into
// a full expression for the expr language.
type ExprBuilder interface {
	Build(cond policies.PolicyCondition) (string, error)
}

type exprBuilder struct {
	builders map[policies.Operator]ExprBuilder
}

func NewExprBuilder() ExprBuilder {
	builders := map[policies.Operator]ExprBuilder{
		// Comparison
		policies.OpEqual:          &ComparisonExprBuilder{},
		policies.OpNotEqual:       &ComparisonExprBuilder{},
		policies.OpLess:           &ComparisonExprBuilder{},
		policies.OpLessOrEqual:    &ComparisonExprBuilder{},
		policies.OpGreater:        &ComparisonExprBuilder{},
		policies.OpGreaterOrEqual: &ComparisonExprBuilder{},

		// Logical
		policies.OpAnd: &LogicalExprBuilder{},
		policies.OpOr:  &LogicalExprBuilder{},
		policies.OpNot: &LogicalExprBuilder{},

		// Set
		policies.OpIn:    &SetExprBuilder{},
		policies.OpNotIn: &SetExprBuilder{},

		// String functions
		policies.OpContains:    &FunctionExprBuilder{},
		policies.OpNotContains: &FunctionExprBuilder{},
		policies.OpStartsWith:  &FunctionExprBuilder{},
		policies.OpEndsWith:    &FunctionExprBuilder{},
		policies.OpMatches:     &FunctionExprBuilder{},

		// Temporal
		policies.OpBefore: &TemporalExprBuilder{},
		policies.OpAfter:  &TemporalExprBuilder{},

		// Arithmetic
		policies.OpMod: &ArithmeticExprBuilder{},
	}

	return &exprBuilder{builders: builders}
}

func (h *exprBuilder) Build(cond policies.PolicyCondition) (string, error) {
	builder, ok := h.builders[cond.Operator]
	if !ok {
		return "", fmt.Errorf("no builder for operator %s", cond.Operator)
	}
	return builder.Build(cond)
}

// LogicalExprBuilder builds logical expressions (and/or/not) for the expr language.
type LogicalExprBuilder struct {
}

// Build builds a logical expression for cond.
func (h *LogicalExprBuilder) Build(cond policies.PolicyCondition) (string, error) {
	if len(cond.Conditions) == 0 {
		return "", fmt.Errorf("logical operator %s requires conditions", cond.Operator)
	}

	var op string
	switch cond.Operator {
	case policies.OpAnd:
		op = "&&"
	case policies.OpOr:
		op = "||"
	case policies.OpNot:
		// Not é unário, espera 1 condição
		if len(cond.Conditions) != 1 {
			return "", fmt.Errorf("operator 'not' requires exactly one condition")
		}
		expr, err := NewExprBuilder().Build(cond.Conditions[0])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("!%s", expr), nil
	default:
		return "", fmt.Errorf("unsupported logical operator: %s", cond.Operator)
	}

	parts := make([]string, 0, len(cond.Conditions))
	for _, c := range cond.Conditions {
		expr, err := NewExprBuilder().Build(c)
		if err != nil {
			return "", err
		}
		parts = append(parts, fmt.Sprintf("(%s)", expr))
	}

	return joinWithOperator(parts, op), nil
}

func joinWithOperator(parts []string, op string) string {
	return stringJoin(parts, " "+op+" ")
}

func stringJoin(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}

// ComparisonExprBuilder builds comparison expressions (==, !=, <, >, etc.).
type ComparisonExprBuilder struct {
}

// Build builds a comparison expression for cond.
func (h *ComparisonExprBuilder) Build(cond policies.PolicyCondition) (string, error) {
	// Exemplo: attribute == value
	// Ajuste para o operador correto
	var op string
	switch cond.Operator {
	case policies.OpEqual:
		op = "=="
	case policies.OpNotEqual:
		op = "!="
	case policies.OpLess:
		op = "<"
	case policies.OpLessOrEqual:
		op = "<="
	case policies.OpGreater:
		op = ">"
	case policies.OpGreaterOrEqual:
		op = ">="
	default:
		return "", fmt.Errorf("unsupported comparison operator: %s", cond.Operator)
	}

	// Supondo que Value pode ser string, número, etc. Ajuste a formatação conforme necessário
	return fmt.Sprintf("%s %s %v", cond.Attribute, op, cond.Value), nil
}

// ArithmeticExprBuilder builds arithmetic expressions (currently mod).
type ArithmeticExprBuilder struct{}

// Build builds an arithmetic expression for cond.
func (h *ArithmeticExprBuilder) Build(cond policies.PolicyCondition) (string, error) {
	if cond.Operator != policies.OpMod {
		return "", fmt.Errorf("unsupported arithmetic operator: %s", cond.Operator)
	}

	return fmt.Sprintf("(%s %% %v)", cond.Attribute, cond.Value), nil
}

// SetExprBuilder builds set membership expressions.
type SetExprBuilder struct{}

// Build builds a set expression for cond.
func (h *SetExprBuilder) Build(cond policies.PolicyCondition) (string, error) {
	var op string
	switch cond.Operator {
	case policies.OpIn:
		op = "in"
	case policies.OpNotIn:
		op = "not in"
	default:
		return "", fmt.Errorf("unsupported set operator: %s", cond.Operator)
	}

	// Value deve ser uma lista/array
	values, ok := cond.Value.([]any)
	if !ok {
		return "", fmt.Errorf("operator %s requires a slice value", cond.Operator)
	}

	// Construir lista de valores
	valStrs := make([]string, len(values))
	for i, v := range values {
		valStrs[i] = fmt.Sprintf("%v", v)
	}

	return fmt.Sprintf("%s %s [%s]", cond.Attribute, op, stringJoin(valStrs, ", ")), nil
}

// FunctionExprBuilder builds string/function expressions (contains, startsWith, etc.).
type FunctionExprBuilder struct{}

// Build builds a function expression for cond.
func (h *FunctionExprBuilder) Build(cond policies.PolicyCondition) (string, error) {
	switch cond.Operator {
	case policies.OpContains:
		return fmt.Sprintf("contains(%s, %v)", cond.Attribute, cond.Value), nil
	case policies.OpNotContains:
		return fmt.Sprintf("!contains(%s, %v)", cond.Attribute, cond.Value), nil
	case policies.OpStartsWith:
		return fmt.Sprintf("startsWith(%s, %v)", cond.Attribute, cond.Value), nil
	case policies.OpEndsWith:
		return fmt.Sprintf("endsWith(%s, %v)", cond.Attribute, cond.Value), nil
	case policies.OpMatches:
		return fmt.Sprintf("matches(%s, %v)", cond.Attribute, cond.Value), nil
	default:
		return "", fmt.Errorf("unsupported function operator: %s", cond.Operator)
	}
}

// TemporalExprBuilder builds temporal comparison expressions (before/after).
type TemporalExprBuilder struct{}

// Build builds a temporal expression for cond.
func (h *TemporalExprBuilder) Build(cond policies.PolicyCondition) (string, error) {
	var op string
	switch cond.Operator {
	case policies.OpBefore:
		op = "<"
	case policies.OpAfter:
		op = ">"
	default:
		return "", fmt.Errorf("unsupported temporal operator: %s", cond.Operator)
	}

	return fmt.Sprintf("%s %s %v", cond.Attribute, op, cond.Value), nil
}
