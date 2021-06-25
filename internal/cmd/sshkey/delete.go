package sshkey

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var deleteCmd = base.DeleteCmd{
	ResourceNameSingular: "SSH Key",
	ShortDescription:     "Delete a SSH Key",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.SSHKey().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.SSHKey().Get(ctx, idOrName)
	},
	Delete: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}) error {
		sshKey := resource.(*hcloud.SSHKey)
		if _, err := client.SSHKey().Delete(ctx, sshKey); err != nil {
			return err
		}
		return nil
	},
}
