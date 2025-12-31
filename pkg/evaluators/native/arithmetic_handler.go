package native

import (
	"fmt"
	"math"

	"github.com/tavaresphil/go-policy-engine/pkg/policies"
	"github.com/tavaresphil/go-policy-engine/pkg/utils"
)

// ArithmeticHandler implements arithmetic operators (currently only mod).
type ArithmeticHandler struct{}

// NewArithmeticHandler constructs an ArithmeticHandler.
func NewArithmeticHandler() OperatorHandler {
	return &ArithmeticHandler{}
}

// Eval implements the modulus check: attribute % value == 0.
// Supports numeric types and numeric strings via utils.AnyToFloat64.
func (h *ArithmeticHandler) Eval(pc policies.PolicyCondition, attr policies.Resolver) (bool, error) {
	if pc.Operator != policies.OpMod {
		return false, fmt.Errorf("unsupported arithmetic operator: %s", pc.Operator)
	}

	attrVal, ok := attr.Resolve(pc.Attribute)
	if !ok {
		return false, fmt.Errorf("missing required attribute: %s", pc.Attribute)
	}

	divisor := pc.Value

	// parse numbers using helper
	af, err := utils.AnyToFloat64(attrVal)
	if err != nil {
		return false, err
	}
	vf, err := utils.AnyToFloat64(divisor)
	if err != nil {
		return false, err
	}

	if vf == 0 {
		return false, fmt.Errorf("modulo by zero")
	}

	res := math.Mod(af, vf)
	// consider floating point precision
	if math.Abs(res) < 1e-9 {
		return true, nil
	}
	return false, nil
}
