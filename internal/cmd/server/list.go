package server

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var serverStatusStrings = []string{
	string(hcloud.ServerStatusInitializing),
	string(hcloud.ServerStatusOff),
	string(hcloud.ServerStatusRunning),
	string(hcloud.ServerStatusStarting),
	string(hcloud.ServerStatusStopping),
	string(hcloud.ServerStatusMigrating),
	string(hcloud.ServerStatusRebuilding),
	string(hcloud.ServerStatusDeleting),
	string(hcloud.ServerStatusUnknown),
}

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Servers",
	JSONKeyGetByName:   "servers",

	DefaultColumns: []string{"id", "name", "status", "ipv4", "ipv6", "private_net", "datacenter", "age"},

	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().StringSlice("status", nil, "Only servers with one of these statuses are displayed")
		_ = cmd.RegisterFlagCompletionFunc("status", cmpl.SuggestCandidates(serverStatusStrings...))
	},

	Fetch: func(s state.State, flags *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		statuses, _ := flags.GetStringSlice("status")

		opts := hcloud.ServerListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		if len(statuses) > 0 {
			for _, status := range statuses {
				if slices.Contains(serverStatusStrings, status) {
					opts.Status = append(opts.Status, hcloud.ServerStatus(status))
				} else {
					return nil, fmt.Errorf("invalid status: %s", status)
				}
			}
		}
		servers, err := s.Client().Server().AllWithOpts(s, opts)

		var resources []interface{}
		for _, r := range servers {
			resources = append(resources, r)
		}
		return resources, err
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.Server{}).
			AddFieldFn("ipv4", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				if server.PublicNet.IPv4.IsUnspecified() {
					return "-"
				}
				return server.PublicNet.IPv4.IP.String()
			})).
			AddFieldFn("ipv6", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				if server.PublicNet.IPv6.IsUnspecified() {
					return "-"
				}
				return server.PublicNet.IPv6.Network.String()
			})).
			AddFieldFn("included_traffic", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				return humanize.IBytes(server.IncludedTraffic)
			})).
			AddFieldFn("ingoing_traffic", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				return humanize.IBytes(server.IngoingTraffic)
			})).
			AddFieldFn("outgoing_traffic", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				return humanize.IBytes(server.OutgoingTraffic)
			})).
			AddFieldFn("datacenter", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				return server.Datacenter.Name
			})).
			AddFieldFn("location", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				return server.Datacenter.Location.Name
			})).
			AddFieldFn("labels", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				return util.LabelsToString(server.Labels)
			})).
			AddFieldFn("type", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				return server.ServerType.Name
			})).
			AddFieldFn("volumes", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				var volumes []string
				for _, volume := range server.Volumes {
					volumeID := strconv.FormatInt(volume.ID, 10)
					volumes = append(volumes, volumeID)
				}
				return strings.Join(volumes, ", ")
			})).
			AddFieldFn("private_net", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				var networks []string
				for _, network := range server.PrivateNet {
					networks = append(networks, fmt.Sprintf("%s (%s)", network.IP.String(), client.Network().Name(network.Network.ID)))
				}
				return util.NA(strings.Join(networks, ", "))
			})).
			AddFieldFn("protection", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				var protection []string
				if server.Protection.Delete {
					protection = append(protection, "delete")
				}
				if server.Protection.Rebuild {
					protection = append(protection, "rebuild")
				}
				return strings.Join(protection, ", ")
			})).
			AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				return util.Datetime(server.Created)
			})).
			AddFieldFn("age", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				return util.Age(server.Created, time.Now())
			})).
			AddFieldFn("placement_group", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				if server.PlacementGroup == nil {
					return "-"
				}
				return server.PlacementGroup.Name
			}))
	},

	Schema: func(resources []interface{}) interface{} {
		serversSchema := make([]schema.Server, 0, len(resources))
		for _, resource := range resources {
			server := resource.(*hcloud.Server)
			serversSchema = append(serversSchema, hcloud.SchemaFromServer(server))
		}
		return serversSchema
	},
}
