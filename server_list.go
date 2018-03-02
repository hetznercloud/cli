package cli

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list [FLAGS]",
		Short:                 "List servers",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerList),
	}
	return cmd
}

func runServerList(cli *CLI, cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetStringArray("output")
	outOpts, err := parseOutputOpts(out)
	if err != nil {
		return err
	}

	servers, err := cli.Client().Server.All(cli.Context)
	if err != nil {
		return err
	}

	cols := []string{"id", "name", "status", "ipv4", "ipv6"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := newTableOutput().
		AddAllowedFields(hcloud.Server{}).
		AddFieldOutputFn("ipv4", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return server.PublicNet.IPv4.IP.String()
		})).
		AddFieldOutputFn("ipv6", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return server.PublicNet.IPv6.Network.String()
		}))

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
