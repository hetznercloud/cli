package floatingip

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd[*hcloud.FloatingIP]{
	ResourceNameSingular: "Floating IP",
	ResourceNamePlural:   "Floating IPs",
	ShortDescription:     "Delete a Floating IP",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.FloatingIP().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.FloatingIP, *hcloud.Response, error) {
		return s.Client().FloatingIP().Get(cmd.Context(), idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, floatingIP *hcloud.FloatingIP) ([]*hcloud.Action, error) {
		_, err := s.Client().FloatingIP().Delete(cmd.Context(), floatingIP)
		return nil, err
	},
}
