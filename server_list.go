package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newServerListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "list",
		Short:            "List servers",
		TraverseChildren: true,
		RunE:             cli.wrap(runServerList),
	}
	cmd.Flags().Bool("no-header", false, "Do not print table header")
	return cmd
}

func runServerList(cli *CLI, cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	noHeader, _ := cmd.Flags().GetBool("no-header")

	servers, err := cli.Client().Server.All(ctx)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	if !noHeader {
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tIPV4")
	}
	for _, server := range servers {
		fmt.Fprintf(w, "%d\t%.50s\t%s\t%s\n", server.ID, server.Name, server.Status,
			server.PublicNet.IPv4.IP)
	}
	w.Flush()

	return nil
}
