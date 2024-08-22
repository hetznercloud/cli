package floatingip

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Floating IP",
	ResourceNamePlural:   "Floating IPs",
	ShortDescription:     "Delete a Floating IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().FloatingIP().Get(s, idOrName)
	},
	Delete: func(s state.State, _ *cobra.Command, resource interface{}) (*hcloud.Action, error) {
		floatingIP := resource.(*hcloud.FloatingIP)
		_, err := s.Client().FloatingIP().Delete(s, floatingIP)
		return nil, err
	},
}
