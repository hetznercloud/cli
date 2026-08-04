package hcapi2

import (
	"context"
	"strconv"
)

type CompletionFunc func(context.Context) ([]string, error)

type LabelCompletionFunc func(context.Context, string) ([]string, error)

func resourceNames[T any](resources []*T, id func(*T) int64, name func(*T) string) []string {
	names := make([]string, len(resources))
	for i, resource := range resources {
		resourceName := name(resource)
		if resourceName == "" {
			resourceName = strconv.FormatInt(id(resource), 10)
		}
		names[i] = resourceName
	}
	return names
}
