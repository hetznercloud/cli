package primaryip

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var updateCmd = base.UpdateCmd{
	ResourceNameSingular: "Primary IP",
	ShortDescription:     "Update a Primary IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PrimaryIP().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.PrimaryIP().Get(ctx, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Primary IP name")
		cmd.Flags().Bool("auto-delete", false, "Delete this Primary IP when the resource it is assigned to is deleted")
	},
	Update: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		primaryIP := resource.(*hcloud.PrimaryIP)
		updOpts := hcloud.PrimaryIPUpdateOpts{
			Name: flags["name"].String(),
		}
		autoDelete, _ := cmd.Flags().GetBool("auto-delete")
		if primaryIP.AutoDelete != autoDelete {
			updOpts.AutoDelete = hcloud.Bool(autoDelete)
		}
		_, _, err := client.PrimaryIP().Update(ctx, primaryIP, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
