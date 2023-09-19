package util

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

func GenerateURLPath(name string) string {
	t := transform.Chain(norm.NFD)
	name, _, _ = transform.String(t, name)

	name = removeCombiningChars(name)

	name = strings.ReplaceAll(name, " ", "-")

	reg, _ := regexp.Compile("[^a-zA-Z0-9-]+")
	name = reg.ReplaceAllString(name, "")

	name = strings.ToLower(name)

	return name
}

func removeCombiningChars(s string) string {
	var result []rune
	for _, r := range s {
		if unicode.Is(unicode.Mn, r) { // Mn: non-spacing marks
			continue
		}
		result = append(result, r)
	}
	return string(result)
}
