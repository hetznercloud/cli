package sshkey

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var labelCmds = base.LabelCmds{
	ResourceNameSingular:   "SSH Key",
	ShortDescriptionAdd:    "Add a label to a SSH Key",
	ShortDescriptionRemove: "Remove a label from a SSH Key",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.SSHKey().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.SSHKey().LabelKeys },
	FetchLabels: func(ctx context.Context, client hcapi2.Client, idOrName string) (map[string]string, int64, error) {
		sshKey, _, err := client.SSHKey().Get(ctx, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if sshKey == nil {
			return nil, 0, fmt.Errorf("ssh key not found: %s", idOrName)
		}
		return sshKey.Labels, sshKey.ID, nil
	},
	SetLabels: func(ctx context.Context, client hcapi2.Client, id int64, labels map[string]string) error {
		opts := hcloud.SSHKeyUpdateOpts{
			Labels: labels,
		}
		_, _, err := client.SSHKey().Update(ctx, &hcloud.SSHKey{ID: id}, opts)
		return err
	},
}
