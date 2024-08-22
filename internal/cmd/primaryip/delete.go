package primaryip

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Primary IP",
	ResourceNamePlural:   "Primary IPs",
	ShortDescription:     "Delete a Primary IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PrimaryIP().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().PrimaryIP().Get(s, idOrName)
	},
	Delete: func(s state.State, _ *cobra.Command, resource interface{}) (*hcloud.Action, error) {
		primaryIP := resource.(*hcloud.PrimaryIP)
		_, err := s.Client().PrimaryIP().Delete(s, primaryIP)
		return nil, err
	},
}
