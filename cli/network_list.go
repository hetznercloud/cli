package cli

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var networkListTableOutput *tableOutput

func init() {
	networkListTableOutput = describeNetworkListTableOutput(nil)
}

func newNetworkListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List networks",
		Long: listLongDescription(
			"Displays a list of networks.",
			networkListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runNetworkList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(networkListTableOutput.Columns()), outputOptionJSON())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runNetworkList(cli *CLI, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	labelSelector, _ := cmd.Flags().GetString("selector")
	opts := hcloud.NetworkListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
			PerPage:       50,
		},
	}
	networks, err := cli.Client().Network.AllWithOpts(cli.Context, opts)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		var networkSchemas []schema.Network
		for _, network := range networks {
			networkSchema := schema.Network{
				ID:         network.ID,
				Name:       network.Name,
				IPRange:    network.IPRange.String(),
				Protection: schema.NetworkProtection{Delete: network.Protection.Delete},
				Created:    network.Created,
				Labels:     network.Labels,
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
		return describeJSON(networkSchemas)
	}

	cols := []string{"id", "name", "ip_range", "servers"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := describeNetworkListTableOutput(cli)
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, network := range networks {
		tw.Write(cols, network)
	}
	tw.Flush()
	return nil
}

func describeNetworkListTableOutput(cli *CLI) *tableOutput {
	return newTableOutput().
		AddAllowedFields(hcloud.Network{}).
		AddFieldOutputFn("servers", fieldOutputFn(func(obj interface{}) string {
			network := obj.(*hcloud.Network)
			serverCount := len(network.Servers)
			if serverCount <= 1 {
				return fmt.Sprintf("%v server", serverCount)
			}
			return fmt.Sprintf("%v servers", serverCount)
		})).
		AddFieldOutputFn("ip_range", fieldOutputFn(func(obj interface{}) string {
			network := obj.(*hcloud.Network)
			return network.IPRange.String()
		})).
		AddFieldOutputFn("labels", fieldOutputFn(func(obj interface{}) string {
			network := obj.(*hcloud.Network)
			return labelsToString(network.Labels)
		})).
		AddFieldOutputFn("protection", fieldOutputFn(func(obj interface{}) string {
			network := obj.(*hcloud.Network)
			var protection []string
			if network.Protection.Delete {
				protection = append(protection, "delete")
			}
			return strings.Join(protection, ", ")
		})).
		AddFieldOutputFn("created", fieldOutputFn(func(obj interface{}) string {
			network := obj.(*hcloud.Network)
			return datetime(network.Created)
		}))
}
