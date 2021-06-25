package floatingip

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var updateCmd = base.UpdateCmd{
	ResourceNameSingular: "Floating IP",
	ShortDescription:     "Update a Floating IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.FloatingIP().Get(ctx, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Floating IP name")
		cmd.Flags().String("description", "", "Floating IP description")
	},
	Update: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		floatingIP := resource.(*hcloud.FloatingIP)
		updOpts := hcloud.FloatingIPUpdateOpts{
			Name:        flags["name"].String(),
			Description: flags["description"].String(),
		}
		_, _, err := client.FloatingIP().Update(ctx, floatingIP, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
