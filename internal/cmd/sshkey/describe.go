package sshkey

import (
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "SSH Key",
	ShortDescription:     "Describe a SSH Key",
	JSONKeyGetByID:       "ssh_key",
	JSONKeyGetByName:     "ssh_keys",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.SSHKey().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (interface{}, interface{}, error) {
		key, _, err := s.Client().SSHKey().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return key, hcloud.SchemaFromSSHKey(key), nil
	},
	PrintText: func(_ state.State, cmd *cobra.Command, resource interface{}) error {
		sshKey := resource.(*hcloud.SSHKey)
		cmd.Printf("ID:\t\t%d\n", sshKey.ID)
		cmd.Printf("Name:\t\t%s\n", sshKey.Name)
		cmd.Printf("Created:\t%s (%s)\n", util.Datetime(sshKey.Created), humanize.Time(sshKey.Created))
		cmd.Printf("Fingerprint:\t%s\n", sshKey.Fingerprint)
		cmd.Printf("Public Key:\n%s\n", strings.TrimSpace(sshKey.PublicKey))
		cmd.Print("Labels:\n")
		if len(sshKey.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range sshKey.Labels {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}

		return nil
	},
}
