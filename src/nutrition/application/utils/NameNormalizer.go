package utils

import (
	"regexp"
	"strings"
)

func NormalizeFoodName(name string) string {
	// 1. Lowercase
	name = strings.ToLower(name)

	// 2. Remove common qualifiers after commas (e.g., "Bananas, raw" -> "bananas")
	if idx := strings.Index(name, ","); idx != -1 {
		name = name[:idx]
	}

	// 3. Remove parts in parentheses
	re := regexp.MustCompile(`\s*\(.*?\)\s*`)
	name = re.ReplaceAllString(name, "")

	// 4. Remove generic words like "raw", "cooked", "dry", "canned", etc.
	genericWords := []string{"raw", "cooked", "dry", "canned", "boiled", "baked", "frozen", "fresh", "prepared"}
	for _, word := range genericWords {
		name = strings.ReplaceAll(name, word, "")
	}

	// 5. Clean extra spaces
	name = strings.TrimSpace(name)
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")

	return name
}
