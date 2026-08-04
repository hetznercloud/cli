package firewall

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd[*hcloud.Firewall]{
	ResourceNameSingular: "Firewall",
	ShortDescription:     "Update a Firewall",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.Firewall().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.Firewall, *hcloud.Response, error) {
		return s.Client().Firewall().Get(cmd.Context(), idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Firewall name")
	},
	Update: func(s state.State, cmd *cobra.Command, firewall *hcloud.Firewall, flags map[string]pflag.Value) error {
		updOpts := hcloud.FirewallUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := s.Client().Firewall().Update(cmd.Context(), firewall, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
