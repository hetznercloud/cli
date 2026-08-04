package primaryip

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd[*hcloud.PrimaryIP]{
	ResourceNameSingular: "Primary IP",
	ResourceNamePlural:   "Primary IPs",
	ShortDescription:     "Delete a Primary IP",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.PrimaryIP().Names(false, false, nil) },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.PrimaryIP, *hcloud.Response, error) {
		return s.Client().PrimaryIP().Get(cmd.Context(), idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, primaryIP *hcloud.PrimaryIP) ([]*hcloud.Action, error) {
		_, err := s.Client().PrimaryIP().Delete(cmd.Context(), primaryIP)
		return nil, err
	},
}
