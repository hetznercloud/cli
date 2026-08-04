package floatingip

import (
	"net"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var SetRDNSCmd = base.SetRdnsCmd[*hcloud.FloatingIP]{
	ResourceNameSingular: "Floating IP",
	ShortDescription:     "Change reverse DNS of a Floating IP",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.FloatingIP().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.FloatingIP, *hcloud.Response, error) {
		return s.Client().FloatingIP().Get(cmd.Context(), idOrName)
	},
	GetDefaultIP: func(floatingIP *hcloud.FloatingIP) net.IP {
		return floatingIP.IP
	},
}
