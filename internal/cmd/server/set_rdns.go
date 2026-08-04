package server

import (
	"net"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var SetRDNSCmd = base.SetRdnsCmd[*hcloud.Server]{
	ResourceNameSingular: "Server",
	ShortDescription:     "Change reverse DNS of a Server",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.Server().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.Server, *hcloud.Response, error) {
		return s.Client().Server().Get(cmd.Context(), idOrName)
	},
	GetDefaultIP: func(server *hcloud.Server) net.IP {
		return server.PublicNet.IPv4.IP
	},
}
