package sshkey

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.SSHKey]{
	ResourceNameSingular:   "SSH Key",
	ShortDescriptionAdd:    "Add a label to a SSH Key",
	ShortDescriptionRemove: "Remove a label from a SSH Key",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.SSHKey().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.SSHKey().LabelKeys },
	Fetch: func(s state.State, idOrName string) (*hcloud.SSHKey, error) {
		sshKey, _, err := s.Client().SSHKey().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if sshKey == nil {
			return nil, fmt.Errorf("ssh key not found: %s", idOrName)
		}
		return sshKey, nil
	},
	SetLabels: func(s state.State, sshKey *hcloud.SSHKey, labels map[string]string) error {
		opts := hcloud.SSHKeyUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().SSHKey().Update(s, sshKey, opts)
		return err
	},
	GetLabels: func(sshKey *hcloud.SSHKey) map[string]string {
		return sshKey.Labels
	},
	GetIDOrName: func(sshKey *hcloud.SSHKey) string {
		return strconv.FormatInt(sshKey.ID, 10)
	},
}
