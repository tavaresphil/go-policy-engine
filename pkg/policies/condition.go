package policies

import "fmt"

// PolicyCondition represents a condition used in a policy. It is either
// a leaf condition (defined by Attribute, Operator and Value) or a logical
// compound condition containing child Conditions (used for operators like
// "and", "or" and "not").
type PolicyCondition struct {
	Attribute  string            `json:"attribute"`
	Operator   Operator          `json:"operator"`
	Value      any               `json:"value"`
	Conditions []PolicyCondition `json:"conditions"`
}

// Validate verifies that the condition is well-formed for the configured
// operator. It checks operator existence, arity for logical operators and
// presence of Attribute/Value for non-logical operators and validates
// child conditions recursively.
func (c PolicyCondition) Validate() error {
	spec, ok := OperatorSpecOf(c.Operator)
	if !ok {
		return fmt.Errorf("unknown operator: %s", c.Operator)
	}

	if spec.Kind == KindLogical {
		n := len(c.Conditions)
		if n < spec.MinArgs || (spec.MaxArgs != -1 && n > spec.MaxArgs) {
			return fmt.Errorf("operator %s expects %d..%d conditions",
				c.Operator, spec.MinArgs, spec.MaxArgs)
		}

		for _, child := range c.Conditions {
			if err := child.Validate(); err != nil {
				return err
			}
		}
		return nil
	}

	if c.Attribute == "" {
		return fmt.Errorf("operator %s requires attribute", c.Operator)
	}

	if c.Value == nil {
		return fmt.Errorf("operator %s requires value", c.Operator)
	}

	return nil
}
