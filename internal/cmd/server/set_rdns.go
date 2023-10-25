package server

import (
	"context"
	"net"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var SetRDNSCmd = base.SetRdnsCmd{
	ResourceNameSingular: "Server",
	ShortDescription:     "Change reverse DNS of a Server",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Server().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Server().Get(ctx, idOrName)
	},
	GetDefaultIP: func(resource interface{}) net.IP {
		server := resource.(*hcloud.Server)
		return server.PublicNet.IPv4.IP
	},
}
