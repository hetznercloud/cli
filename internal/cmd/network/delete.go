package network

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd[*hcloud.Network]{
	ResourceNameSingular: "Network",
	ResourceNamePlural:   "Networks",
	ShortDescription:     "Delete a network",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.Network().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.Network, *hcloud.Response, error) {
		return s.Client().Network().Get(cmd.Context(), idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, network *hcloud.Network) ([]*hcloud.Action, error) {
		_, err := s.Client().Network().Delete(cmd.Context(), network)
		return nil, err
	},
}
