package storagebox

import (
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/pflag"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Storage Boxes",
	JSONKeyGetByName:   "storage_boxes",
	DefaultColumns:     []string{"id", "name", "type", "username", "size", "location"},
	Fetch: func(s state.State, set *pflag.FlagSet, opts hcloud.ListOpts, strings []string) ([]interface{}, error) {
		storageBoxes, err := s.Client().StorageBox().All(s)

		var resources []interface{}
		for _, r := range storageBoxes {
			resources = append(resources, r)
		}
		return resources, err
	},
	OutputTable: func(t *output.Table, client hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.StorageBox{}).
			AddFieldFn("type", output.FieldFn(func(obj interface{}) string {
				storageBox := obj.(*hcloud.StorageBox)
				return storageBox.StorageBoxType.Name
			})).
			AddFieldFn("username", output.FieldFn(func(obj interface{}) string {
				storageBox := obj.(*hcloud.StorageBox)
				if storageBox.Username == nil {
					return "-"
				}
				return *storageBox.Username
			})).
			AddFieldFn("location", output.FieldFn(func(obj interface{}) string {
				storageBox := obj.(*hcloud.StorageBox)
				return storageBox.Location.Name
			})).
			AddFieldFn("size", output.FieldFn(func(obj interface{}) string {
				storageBox := obj.(*hcloud.StorageBox)
				return strconv.FormatUint(storageBox.Stats.Size, 10)
			}))
	},
}
