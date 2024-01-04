package sshkey

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd{
	ResourceNameSingular: "SSH Key",
	ShortDescription:     "Delete a SSH Key",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.SSHKey().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().SSHKey().Get(s, idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		sshKey := resource.(*hcloud.SSHKey)
		if _, err := s.Client().SSHKey().Delete(s, sshKey); err != nil {
			return err
		}
		return nil
	},
}
