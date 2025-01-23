package sshkey

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "SSH Key",
	ShortDescriptionAdd:    "Add a label to a SSH Key",
	ShortDescriptionRemove: "Remove a label from a SSH Key",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.SSHKey().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.SSHKey().LabelKeys },
	Fetch: func(s state.State, idOrName string) (any, error) {
		sshKey, _, err := s.Client().SSHKey().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if sshKey == nil {
			return nil, fmt.Errorf("ssh key not found: %s", idOrName)
		}
		return sshKey, nil
	},
	SetLabels: func(s state.State, resource any, labels map[string]string) error {
		sshKey := resource.(*hcloud.SSHKey)
		opts := hcloud.SSHKeyUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().SSHKey().Update(s, sshKey, opts)
		return err
	},
	GetLabels: func(resource any) map[string]string {
		sshKey := resource.(*hcloud.SSHKey)
		return sshKey.Labels
	},
	GetIDOrName: func(resource any) string {
		sshKey := resource.(*hcloud.SSHKey)
		return strconv.FormatInt(sshKey.ID, 10)
	},
}
