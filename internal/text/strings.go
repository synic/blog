package text

import (
	"regexp"
	"strings"
)

var (
	nonAlphaNumRe = regexp.MustCompile(`[^a-z0-9-]`)
	multiHyphenRe = regexp.MustCompile(`-+`)
)

func Slugify(title string) string {
	// Create a proper slug from title:
	// 1. Convert to lowercase
	// 2. Replace all spaces and underscores with hyphens
	// 3. Remove all non-alphanumeric characters except hyphens
	// 4. Remove multiple consecutive hyphens
	// 5. Trim hyphens from start/end
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = nonAlphaNumRe.ReplaceAllString(slug, "")
	slug = multiHyphenRe.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}
