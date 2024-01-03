package floatingip

import (
	"net"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var SetRDNSCmd = base.SetRdnsCmd{
	ResourceNameSingular: "Floating IP",
	ShortDescription:     "Change reverse DNS of a Floating IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.FloatingIP().Get(s, idOrName)
	},
	GetDefaultIP: func(resource interface{}) net.IP {
		floatingIP := resource.(*hcloud.FloatingIP)
		return floatingIP.IP
	},
}
