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

var ListCmd = base.ListCmd{
	ResourceNamePlural: "SSH keys",
	JSONKeyGetByName:   "ssh_keys",
	DefaultColumns:     []string{"id", "name", "fingerprint", "age"},
	SortOption:         config.OptionSortSSHKey,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.SSHKeyListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		sshKeys, err := s.Client().SSHKey().AllWithOpts(s, opts)

		var resources []interface{}
		for _, n := range sshKeys {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(t *output.Table, _ hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.SSHKey{}).
			AddFieldFn("labels", output.FieldFn(func(obj interface{}) string {
				sshKey := obj.(*hcloud.SSHKey)
				return util.LabelsToString(sshKey.Labels)
			})).
			AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
				sshKey := obj.(*hcloud.SSHKey)
				return util.Datetime(sshKey.Created)
			})).
			AddFieldFn("age", output.FieldFn(func(obj interface{}) string {
				sshKey := obj.(*hcloud.SSHKey)
				return util.Age(sshKey.Created, time.Now())
			}))
	},

	Schema: func(resources []interface{}) interface{} {
		sshKeySchemas := make([]schema.SSHKey, 0, len(resources))
		for _, resource := range resources {
			sshKey := resource.(*hcloud.SSHKey)
			sshKeySchemas = append(sshKeySchemas, hcloud.SchemaFromSSHKey(sshKey))
		}
		return sshKeySchemas
	},
}
