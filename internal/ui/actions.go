package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func ActionMessage(action *hcloud.Action) string {
	return fmt.Sprintf("Waiting for %s to complete", color.New(color.Bold).Sprint(action.Command))
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
