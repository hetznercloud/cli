package sshkey

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

var listCmd = base.ListCmd{
	ResourceNamePlural: "ssh keys",
	DefaultColumns:     []string{"id", "name", "fingerprint"},

	Fetch: func(ctx context.Context, client hcapi2.Client, listOpts hcloud.ListOpts) ([]interface{}, error) {
		sshKeys, _, err := client.SSHKey().List(ctx, hcloud.SSHKeyListOpts{ListOpts: listOpts})

		var resources []interface{}
		for _, n := range sshKeys {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(_ hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.SSHKey{}).
			AddFieldFn("labels", output.FieldFn(func(obj interface{}) string {
				sshKey := obj.(*hcloud.SSHKey)
				return util.LabelsToString(sshKey.Labels)
			})).
			AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
				sshKey := obj.(*hcloud.SSHKey)
				return util.Datetime(sshKey.Created)
			}))
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var sshKeySchemas []schema.SSHKey
		for _, resource := range resources {
			sshKey := resource.(*hcloud.SSHKey)
			sshKeySchema := schema.SSHKey{
				ID:          sshKey.ID,
				Name:        sshKey.Name,
				Fingerprint: sshKey.Fingerprint,
				PublicKey:   sshKey.PublicKey,
				Labels:      sshKey.Labels,
				Created:     sshKey.Created,
			}
			sshKeySchemas = append(sshKeySchemas, sshKeySchema)
		}
		return sshKeySchemas
	},
}
