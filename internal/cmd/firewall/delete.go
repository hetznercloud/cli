package firewall

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd[*hcloud.Firewall]{
	ResourceNameSingular: "Firewall",
	ResourceNamePlural:   "Firewalls",
	ShortDescription:     "Delete a Firewall",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.Firewall().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.Firewall, *hcloud.Response, error) {
		return s.Client().Firewall().Get(cmd.Context(), idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, firewall *hcloud.Firewall) ([]*hcloud.Action, error) {
		_, err := s.Client().Firewall().Delete(cmd.Context(), firewall)
		return nil, err
	},
}
