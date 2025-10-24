package sshkey

import (
	"time"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = &base.ListCmd[*hcloud.SSHKey, schema.SSHKey]{
	ResourceNamePlural: "SSH Keys",
	JSONKeyGetByName:   "ssh_keys",
	DefaultColumns:     []string{"id", "name", "fingerprint", "age"},
	SortOption:         config.OptionSortSSHKey,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.SSHKey, error) {
		opts := hcloud.SSHKeyListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().SSHKey().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.SSHKey], _ hcapi2.Client) {
		t.
			AddAllowedFields(&hcloud.SSHKey{}).
			AddFieldFn("labels", func(sshKey *hcloud.SSHKey) string {
				return util.LabelsToString(sshKey.Labels)
			}).
			AddFieldFn("created", func(sshKey *hcloud.SSHKey) string {
				return util.Datetime(sshKey.Created)
			}).
			AddFieldFn("age", func(sshKey *hcloud.SSHKey) string {
				return util.Age(sshKey.Created, time.Now())
			})
	},

	Schema: hcloud.SchemaFromSSHKey,
}
