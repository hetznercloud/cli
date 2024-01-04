package placementgroup

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd{
	ResourceNameSingular: "placement group",
	ShortDescription:     "Update a placement group",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PlacementGroup().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().PlacementGroup().Get(s, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Placement group name")
	},
	Update: func(s state.State, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		placementGroup := resource.(*hcloud.PlacementGroup)
		updOpts := hcloud.PlacementGroupUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := s.Client().PlacementGroup().Update(s, placementGroup, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
