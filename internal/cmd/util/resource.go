package util

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ResourceState struct {
	IDOrName string
	Error    error
}

// ResourceDisplayName is a fallback formatting for resources with an unknown type.
// Example: foo_bar 123 -> Foo Bar 123
func ResourceDisplayName(resourceType string, id int64) string {
	return fmt.Sprintf("%s %d", cases.Title(language.English).String(strings.ReplaceAll(resourceType, "_", " ")), id)
}
