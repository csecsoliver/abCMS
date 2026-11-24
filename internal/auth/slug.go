package auth

import (
	"regexp"
	"strings"
)

var nonAlnumRun = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify normalizes a username or label into a lowercase hyphenated slug,
// matching the behavior of Lapis's lapis.util.slugify for ASCII usernames.
func Slugify(input string) string {
	normalized := strings.ToLower(strings.TrimSpace(input))
	normalized = nonAlnumRun.ReplaceAllString(normalized, "-")
	return strings.Trim(normalized, "-")
}
