package firewall

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

var updateCmd = base.UpdateCmd{
	ResourceNameSingular: "Firewall",
	ShortDescription:     "Update a firewall",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Firewall().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Firewall().Get(ctx, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Firewall name")
	},
	Update: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		firewall := resource.(*hcloud.Firewall)
		updOpts := hcloud.FirewallUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := client.Firewall().Update(ctx, firewall, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
