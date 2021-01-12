package network

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var listTableOutput *output.Table

func init() {
	listTableOutput = describelistTableOutput(nil)
}

func newListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List networks",
		Long: util.ListLongDescription(
			"Displays a list of networks.",
			listTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(listTableOutput.Columns()), output.OptionJSON())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

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
		return util.DescribeJSON(networkSchemas)
	}

	cols := []string{"id", "name", "ip_range", "servers"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := describelistTableOutput(cli)
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

func describelistTableOutput(cli *state.State) *output.Table {
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
		}))
}
