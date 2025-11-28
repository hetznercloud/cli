package primaryip

import (
	"net"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var SetRDNSCmd = base.SetRdnsCmd[*hcloud.PrimaryIP]{
	ResourceNameSingular: "Primary IP",
	ShortDescription:     "Change reverse DNS of a Primary IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PrimaryIP().Names(false, false, nil) },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.PrimaryIP, *hcloud.Response, error) {
		return s.Client().PrimaryIP().Get(s, idOrName)
	},
	GetDefaultIP: func(primaryIP *hcloud.PrimaryIP) net.IP {
		return primaryIP.IP
	},
}
