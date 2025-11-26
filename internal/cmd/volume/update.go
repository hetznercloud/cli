package volume

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd[*hcloud.Volume]{
	ResourceNameSingular: "Volume",
	ShortDescription:     "Update a Volume",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Volume().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Volume, *hcloud.Response, error) {
		return s.Client().Volume().Get(s, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Volume name")
	},
	Update: func(s state.State, _ *cobra.Command, volume *hcloud.Volume, flags map[string]pflag.Value) error {
		updOpts := hcloud.VolumeUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := s.Client().Volume().Update(s, volume, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
