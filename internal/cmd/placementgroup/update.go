package placementgroup

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd[*hcloud.PlacementGroup]{
	ResourceNameSingular: "Placement Group",
	ShortDescription:     "Update a Placement Group",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.PlacementGroup().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.PlacementGroup, *hcloud.Response, error) {
		return s.Client().PlacementGroup().Get(cmd.Context(), idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Placement Group name")
	},
	Update: func(s state.State, cmd *cobra.Command, placementGroup *hcloud.PlacementGroup, flags map[string]pflag.Value) error {
		updOpts := hcloud.PlacementGroupUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := s.Client().PlacementGroup().Update(cmd.Context(), placementGroup, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
