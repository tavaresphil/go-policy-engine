package expr

import (
	"encoding/json"
	"fmt"
)

func Literal(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("invalid literal: %w", err)
	}
	return string(b), nil
}
