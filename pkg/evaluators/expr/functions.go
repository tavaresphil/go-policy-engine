package expr

import "strings"

func Functions() map[string]any {
	return map[string]any{
		"contains": func(s, substr string) bool {
			return strings.Contains(s, substr)
		},
		"startsWith": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"endsWith": func(s, suffix string) bool {
			return strings.HasSuffix(s, suffix)
		},
	}
}
