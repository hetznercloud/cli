package server

import (
	"context"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "servers",

	DefaultColumns: []string{"id", "name", "status", "ipv4", "ipv6", "datacenter"},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts) ([]interface{}, error) {
		servers, err := client.Server().AllWithOpts(ctx, hcloud.ServerListOpts{ListOpts: listOpts})

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
				return server.PublicNet.IPv4.IP.String()
			})).
			AddFieldFn("ipv6", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
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
					volumeID := strconv.Itoa(volume.ID)
					volumes = append(volumes, volumeID)
				}
				return strings.Join(volumes, ", ")
			})).
			AddFieldFn("private_net", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				var networks []string
				for _, network := range server.PrivateNet {
					networks = append(networks, client.Network().Name(network.Network.ID))
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
			AddFieldFn("placement_group", output.FieldFn(func(obj interface{}) string {
				server := obj.(*hcloud.Server)
				if server.PlacementGroup == nil {
					return "-"
				}
				return server.PlacementGroup.Name
			}))
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var serversSchema []schema.Server
		for _, resource := range resources {
			server := resource.(*hcloud.Server)

			serverSchema := schema.Server{
				ID:         server.ID,
				Name:       server.Name,
				Status:     string(server.Status),
				Created:    server.Created,
				Datacenter: util.DatacenterToSchema(*server.Datacenter),
				ServerType: util.ServerTypeToSchema(*server.ServerType),
				PublicNet: schema.ServerPublicNet{
					IPv4: schema.ServerPublicNetIPv4{
						IP:      server.PublicNet.IPv4.IP.String(),
						Blocked: server.PublicNet.IPv4.Blocked,
						DNSPtr:  server.PublicNet.IPv4.DNSPtr,
					},
					IPv6: schema.ServerPublicNetIPv6{
						IP:      server.PublicNet.IPv6.IP.String(),
						Blocked: server.PublicNet.IPv6.Blocked,
					},
				},
				RescueEnabled:   server.RescueEnabled,
				BackupWindow:    hcloud.String(server.BackupWindow),
				OutgoingTraffic: &server.OutgoingTraffic,
				IngoingTraffic:  &server.IngoingTraffic,
				IncludedTraffic: server.IncludedTraffic,
				Protection: schema.ServerProtection{
					Delete:  server.Protection.Delete,
					Rebuild: server.Protection.Rebuild,
				},
				Labels:          server.Labels,
				PrimaryDiskSize: server.PrimaryDiskSize,
			}
			if server.Image != nil {
				serverImage := util.ImageToSchema(*server.Image)
				serverSchema.Image = &serverImage
			}
			if server.ISO != nil {
				serverISO := util.ISOToSchema(*server.ISO)
				serverSchema.ISO = &serverISO
			}
			for ip, dnsPTR := range server.PublicNet.IPv6.DNSPtr {
				serverSchema.PublicNet.IPv6.DNSPtr = append(serverSchema.PublicNet.IPv6.DNSPtr, schema.ServerPublicNetIPv6DNSPtr{
					IP:     ip,
					DNSPtr: dnsPTR,
				})
			}
			for _, floatingIP := range server.PublicNet.FloatingIPs {
				serverSchema.PublicNet.FloatingIPs = append(serverSchema.PublicNet.FloatingIPs, floatingIP.ID)
			}
			for _, volume := range server.Volumes {
				serverSchema.Volumes = append(serverSchema.Volumes, volume.ID)
			}
			for _, privateNet := range server.PrivateNet {
				privateNetSchema := schema.ServerPrivateNet{
					Network:    privateNet.Network.ID,
					IP:         privateNet.IP.String(),
					MACAddress: privateNet.MACAddress,
				}
				for _, aliasIP := range privateNet.Aliases {
					privateNetSchema.AliasIPs = append(privateNetSchema.AliasIPs, aliasIP.String())
				}
				serverSchema.PrivateNet = append(serverSchema.PrivateNet, privateNetSchema)
			}
			serversSchema = append(serversSchema, serverSchema)
		}
		return serversSchema
	},
}
