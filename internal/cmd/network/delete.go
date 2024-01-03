package network

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Network",
	ShortDescription:     "Delete a network",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Network().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Network().Get(s, idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		network := resource.(*hcloud.Network)
		if _, err := s.Network().Delete(s, network); err != nil {
			return err
		}
		return nil
	},
}
