package floatingip

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd{
	ResourceNameSingular: "Floating IP",
	ShortDescription:     "Update a Floating IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().FloatingIP().Get(s, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Floating IP name")
		cmd.Flags().String("description", "", "Floating IP description")
	},
	Update: func(s state.State, _ *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		floatingIP := resource.(*hcloud.FloatingIP)
		updOpts := hcloud.FloatingIPUpdateOpts{
			Name:        flags["name"].String(),
			Description: flags["description"].String(),
		}
		_, _, err := s.Client().FloatingIP().Update(s, floatingIP, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
