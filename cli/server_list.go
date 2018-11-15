package cli

import (
	"strconv"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var serverListTableOutput *tableOutput

func init() {
	serverListTableOutput = newTableOutput().
		AddAllowedFields(hcloud.Server{}).
		AddFieldOutputFn("ipv4", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return server.PublicNet.IPv4.IP.String()
		})).
		AddFieldOutputFn("ipv6", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return server.PublicNet.IPv6.Network.String()
		})).
		AddFieldOutputFn("datacenter", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return server.Datacenter.Name
		})).
		AddFieldOutputFn("location", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return server.Datacenter.Location.Name
		})).
		AddFieldOutputFn("labels", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return labelsToString(server.Labels)
		})).
		AddFieldOutputFn("type", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return server.ServerType.Name
		})).
		AddFieldOutputFn("volumes", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			var volumes []string
			for _, volume := range server.Volumes {
				volumeID := strconv.Itoa(volume.ID)
				volumes = append(volumes, volumeID)
			}
			return strings.Join(volumes, ", ")
		})).
		AddFieldOutputFn("protection", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			var protection []string
			if server.Protection.Delete {
				protection = append(protection, "delete")
			}
			if server.Protection.Rebuild {
				protection = append(protection, "rebuild")
			}
			return strings.Join(protection, ", ")
		}))
}

func newServerListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List servers",
		Long: listLongDescription(
			"Displays a list of servers.",
			serverListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerList),
	}
	addListOutputFlag(cmd, serverListTableOutput.Columns())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runServerList(cli *CLI, cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetStringArray("output")
	outOpts, err := parseOutputOpts(out)
	if err != nil {
		return err
	}

	labelSelector, _ := cmd.Flags().GetString("selector")
	opts := hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
			PerPage:       50,
		},
	}
	servers, err := cli.Client().Server.AllWithOpts(cli.Context, opts)
	if err != nil {
		return err
	}

	cols := []string{"id", "name", "status", "ipv4", "ipv6", "datacenter"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := serverListTableOutput
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, server := range servers {
		tw.Write(cols, server)
	}
	tw.Flush()
	return nil
}
