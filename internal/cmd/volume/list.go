package volume

import (
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Volumes",
	JSONKeyGetByName:   "volumes",
	DefaultColumns:     []string{"id", "name", "size", "server", "location", "age"},

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.VolumeListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		volumes, err := s.Client().Volume().AllWithOpts(s, opts)

		var resources []interface{}
		for _, n := range volumes {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.Volume{}).
			AddFieldFn("server", output.FieldFn(func(obj interface{}) string {
				volume := obj.(*hcloud.Volume)
				var server string
				if volume.Server != nil {
					return client.Server().ServerName(volume.Server.ID)
				}
				return util.NA(server)
			})).
			AddFieldFn("size", output.FieldFn(func(obj interface{}) string {
				volume := obj.(*hcloud.Volume)
				return humanize.Bytes(uint64(volume.Size * humanize.GByte))
			})).
			AddFieldFn("location", output.FieldFn(func(obj interface{}) string {
				volume := obj.(*hcloud.Volume)
				return volume.Location.Name
			})).
			AddFieldFn("protection", output.FieldFn(func(obj interface{}) string {
				volume := obj.(*hcloud.Volume)
				var protection []string
				if volume.Protection.Delete {
					protection = append(protection, "delete")
				}
				return strings.Join(protection, ", ")
			})).
			AddFieldFn("labels", output.FieldFn(func(obj interface{}) string {
				volume := obj.(*hcloud.Volume)
				return util.LabelsToString(volume.Labels)
			})).
			AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
				volume := obj.(*hcloud.Volume)
				return util.Datetime(volume.Created)
			})).
			AddFieldFn("age", output.FieldFn(func(obj interface{}) string {
				volume := obj.(*hcloud.Volume)
				return util.Age(volume.Created, time.Now())
			}))
	},

	Schema: func(resources []interface{}) interface{} {
		volumesSchema := make([]schema.Volume, 0, len(resources))
		for _, resource := range resources {
			volume := resource.(*hcloud.Volume)
			volumesSchema = append(volumesSchema, hcloud.SchemaFromVolume(volume))
		}
		return volumesSchema
	},
}
