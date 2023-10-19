package sshkey

import (
	"context"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var describeCmd = base.DescribeCmd{
	ResourceNameSingular: "SSH Key",
	ShortDescription:     "Describe a SSH Key",
	JSONKeyGetByID:       "ssh_key",
	JSONKeyGetByName:     "ssh_keys",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.SSHKey().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.SSHKey().Get(ctx, idOrName)
	},
	PrintText: func(_ context.Context, _ hcapi2.Client, _ *cobra.Command, resource interface{}) error {
		sshKey := resource.(*hcloud.SSHKey)
		fmt.Printf("ID:\t\t%d\n", sshKey.ID)
		fmt.Printf("Name:\t\t%s\n", sshKey.Name)
		fmt.Printf("Created:\t%s (%s)\n", util.Datetime(sshKey.Created), humanize.Time(sshKey.Created))
		fmt.Printf("Fingerprint:\t%s\n", sshKey.Fingerprint)
		fmt.Printf("Public Key:\n%s\n", strings.TrimSpace(sshKey.PublicKey))
		fmt.Print("Labels:\n")
		if len(sshKey.Labels) == 0 {
			fmt.Print("  No labels\n")
		} else {
			for key, value := range sshKey.Labels {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}

		return nil
	},
}
