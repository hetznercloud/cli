package sshkey

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd[*hcloud.SSHKey]{
	ResourceNameSingular: "SSH Key",
	ShortDescription:     "Update an SSH Key",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.SSHKey().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.SSHKey, *hcloud.Response, error) {
		return s.Client().SSHKey().Get(s, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "SSH Key name")
	},
	Update: func(s state.State, _ *cobra.Command, sshKey *hcloud.SSHKey, flags map[string]pflag.Value) error {
		updOpts := hcloud.SSHKeyUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := s.Client().SSHKey().Update(s, sshKey, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
