package util

import (
	"fmt"

	"github.com/dustin/go-humanize"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func DescribeDeprecation(resource hcloud.Deprecatable) string {
	if !resource.IsDeprecated() {
		return ""
	}

	info := "Deprecation:\t\n"
	info += fmt.Sprintf("  Announced:\t%s (%s)\n", Datetime(resource.DeprecationAnnounced()), humanize.Time(resource.DeprecationAnnounced()))
	info += fmt.Sprintf("  Unavailable After:\t%s (%s)\n", Datetime(resource.UnavailableAfter()), humanize.Time(resource.UnavailableAfter()))

	return info
}
