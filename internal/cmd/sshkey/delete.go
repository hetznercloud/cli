package sshkey

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd[*hcloud.SSHKey]{
	ResourceNameSingular: "SSH Key",
	ResourceNamePlural:   "SSH Keys",
	ShortDescription:     "Delete an SSH Key",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.SSHKey().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.SSHKey, *hcloud.Response, error) {
		return s.Client().SSHKey().Get(cmd.Context(), idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, sshKey *hcloud.SSHKey) ([]*hcloud.Action, error) {
		_, err := s.Client().SSHKey().Delete(cmd.Context(), sshKey)
		return nil, err
	},
}
