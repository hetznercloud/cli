package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func ActionMessage(action *hcloud.Action) string {
	return fmt.Sprintf("Waiting for %s", color.New(color.Bold).Sprint(action.Command))
}

// FakeActionMessage returns the initial value with a unused color to grow the string
// size.
//
// Because the [ActionMessage] function adds 1 color to the returned string. We add the
// same amount of colors to the [FakeActionMessage], to make sure the padding is
// correct.
func FakeActionMessage(value string) string {
	return color.New(color.Bold).Sprint("") + value
}

func ActionResourcesMessage(resources ...*hcloud.ActionResource) string {
	if len(resources) == 0 {
		return ""
	}

	items := make([]string, 0, len(resources))
	for _, resource := range resources {
		items = append(items, fmt.Sprintf("%s: %d", resource.Type, resource.ID))
	}

	return fmt.Sprintf("(%v)", strings.Join(items, ", "))
}
