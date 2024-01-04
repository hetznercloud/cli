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
	ShortDescription:     "Delete a Floating IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().FloatingIP().Get(s, idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		floatingIP := resource.(*hcloud.FloatingIP)
		if _, err := s.Client().FloatingIP().Delete(s, floatingIP); err != nil {
			return err
		}
		return nil
	},
}
