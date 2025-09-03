package storagebox

import (
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Storage Boxes",
	JSONKeyGetByName:   "storage_boxes",
	DefaultColumns:     []string{"id", "name", "username", "server", "type", "size", "location", "age"},
	Fetch: func(s state.State, set *pflag.FlagSet, opts hcloud.ListOpts, strings []string) ([]interface{}, error) {
		listOpts := hcloud.StorageBoxListOpts{ListOpts: opts}
		storageBoxes, err := s.Client().StorageBox().AllWithOpts(s, listOpts)

		var resources []interface{}
		for _, r := range storageBoxes {
			resources = append(resources, r)
		}
		return resources, err
	},
	OutputTable: func(t *output.Table, client hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.StorageBox{}).
			AddFieldFn("type", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return storageBox.StorageBoxType.Name
			}).
			AddFieldFn("username", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				if storageBox.Username == nil {
					return "-"
				}
				return *storageBox.Username
			}).
			AddFieldFn("server", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				if storageBox.Server == nil {
					return "-"
				}
				return *storageBox.Server
			}).
			AddFieldFn("system", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				if storageBox.System == nil {
					return "-"
				}
				return *storageBox.System
			}).
			AddFieldFn("location", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return storageBox.Location.Name
			}).
			AddFieldFn("size", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				if storageBox.Stats == nil {
					return "-"
				}
				return humanize.IBytes(storageBox.Stats.Size)
			}).
			AddFieldFn("labels", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return util.LabelsToString(storageBox.Labels)
			}).
			AddFieldFn("created", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return util.Datetime(storageBox.Created)
			}).
			AddFieldFn("age", func(obj any) string {
				storageBox := obj.(*hcloud.StorageBox)
				return util.Age(storageBox.Created, time.Now())
			})
	},
}
