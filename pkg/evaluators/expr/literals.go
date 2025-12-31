package expr

import (
	"encoding/json"
	"fmt"
)

// Literal returns a JSON literal representation of v suitable for embedding
// in an expression built by the expr engine.
func Literal(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("invalid literal: %w", err)
	}
	return string(b), nil
}
