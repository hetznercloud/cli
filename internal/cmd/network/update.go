package network

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var updateCmd = base.UpdateCmd{
	ResourceNameSingular: "Network",
	ShortDescription:     "Update a Network.\n\nTo enable or disable exposing routes to the vSwitch connection you can use the subcommand \"hcloud network expose-routes-to-vswitch\".",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Network().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Network().Get(ctx, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Network name")
	},
	Update: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		floatingIP := resource.(*hcloud.Network)
		updOpts := hcloud.NetworkUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := client.Network().Update(ctx, floatingIP, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
