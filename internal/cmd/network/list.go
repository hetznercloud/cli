package network

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Networks",
	DefaultColumns:     []string{"id", "name", "ip_range", "servers", "age"},

	Fetch: func(ctx context.Context, client hcapi2.Client, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.NetworkListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		networks, err := client.Network().AllWithOpts(ctx, opts)

		var resources []interface{}
		for _, n := range networks {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(_ hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.Network{}).
			AddFieldFn("servers", output.FieldFn(func(obj interface{}) string {
				network := obj.(*hcloud.Network)
				serverCount := len(network.Servers)
				if serverCount <= 1 {
					return fmt.Sprintf("%v server", serverCount)
				}
				return fmt.Sprintf("%v servers", serverCount)
			})).
			AddFieldFn("ip_range", output.FieldFn(func(obj interface{}) string {
				network := obj.(*hcloud.Network)
				return network.IPRange.String()
			})).
			AddFieldFn("labels", output.FieldFn(func(obj interface{}) string {
				network := obj.(*hcloud.Network)
				return util.LabelsToString(network.Labels)
			})).
			AddFieldFn("protection", output.FieldFn(func(obj interface{}) string {
				network := obj.(*hcloud.Network)
				var protection []string
				if network.Protection.Delete {
					protection = append(protection, "delete")
				}
				return strings.Join(protection, ", ")
			})).
			AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
				network := obj.(*hcloud.Network)
				return util.Datetime(network.Created)
			})).
			AddFieldFn("age", output.FieldFn(func(obj interface{}) string {
				network := obj.(*hcloud.Network)
				return util.Age(network.Created, time.Now())
			}))
	},
	JSONSchema: func(resources []interface{}) interface{} {
		networkSchemas := make([]schema.Network, 0, len(resources))
		for _, resource := range resources {
			network := resource.(*hcloud.Network)

			networkSchema := schema.Network{
				ID:                    network.ID,
				Name:                  network.Name,
				IPRange:               network.IPRange.String(),
				Protection:            schema.NetworkProtection{Delete: network.Protection.Delete},
				Created:               network.Created,
				Labels:                network.Labels,
				ExposeRoutesToVSwitch: network.ExposeRoutesToVSwitch,
			}
			for _, subnet := range network.Subnets {
				networkSchema.Subnets = append(networkSchema.Subnets, schema.NetworkSubnet{
					Type:        string(subnet.Type),
					IPRange:     subnet.IPRange.String(),
					NetworkZone: string(subnet.NetworkZone),
					Gateway:     subnet.Gateway.String(),
				})
			}
			for _, route := range network.Routes {
				networkSchema.Routes = append(networkSchema.Routes, schema.NetworkRoute{
					Destination: route.Destination.String(),
					Gateway:     route.Gateway.String(),
				})
			}
			for _, server := range network.Servers {
				networkSchema.Servers = append(networkSchema.Servers, server.ID)
			}
			networkSchemas = append(networkSchemas, networkSchema)
		}
		return networkSchemas
	},
}
