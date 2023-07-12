package sshkey

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var updateCmd = base.UpdateCmd{
	ResourceNameSingular: "SSHKey",
	ShortDescription:     "Update a SSHKey",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.SSHKey().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.SSHKey().Get(ctx, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "SSH Key name")
	},
	Update: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		floatingIP := resource.(*hcloud.SSHKey)
		updOpts := hcloud.SSHKeyUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := client.SSHKey().Update(ctx, floatingIP, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
