package volume

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var updateCmd = base.UpdateCmd{
	ResourceNameSingular: "Volume",
	ShortDescription:     "Update a Volume",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Volume().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Volume().Get(ctx, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Volume name")
	},
	Update: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		floatingIP := resource.(*hcloud.Volume)
		updOpts := hcloud.VolumeUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := client.Volume().Update(ctx, floatingIP, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
