package firewall

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd{
	ResourceNameSingular: "Firewall",
	ShortDescription:     "Update a firewall",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Firewall().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().Firewall().Get(s, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Firewall name")
	},
	Update: func(s state.State, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		firewall := resource.(*hcloud.Firewall)
		updOpts := hcloud.FirewallUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := s.Client().Firewall().Update(s, firewall, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
