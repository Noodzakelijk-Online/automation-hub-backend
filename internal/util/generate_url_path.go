package util

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

func GenerateURLPath(name string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	name, _, _ = transform.String(t, name)

	name = strings.ReplaceAll(name, " ", "-")

	reg, _ := regexp.Compile("[^a-zA-Z0-9-]+")
	name = reg.ReplaceAllString(name, "")

	name = strings.ToLower(name)

	return name
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}
