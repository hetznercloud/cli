package floatingip

import (
	"context"
	"net"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

var setRDNSCmd = base.SetRdnsCmd{
	ResourceNameSingular: "Floating IP",
	ShortDescription:     "Change reverse DNS of a Floating IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.FloatingIP().Get(ctx, idOrName)
	},
	GetDefaultIP: func(resource interface{}) net.IP {
		floatingIP := resource.(*hcloud.FloatingIP)
		return floatingIP.IP
	},
}
