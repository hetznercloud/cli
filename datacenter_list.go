package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newDatacenterListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list [FLAGS]",
		Short:                 "List datacenters",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runDatacenterList),
	}
	return cmd
}

func runDatacenterList(cli *CLI, cmd *cobra.Command, args []string) error {
	datacenters, err := cli.Client().Datacenter.All(cli.Context)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION\tLOCATION")
	for _, datacenter := range datacenters {
		fmt.Fprintf(w, "%d\t%.50s\t%.50s\t%s\n",
			datacenter.ID,
			datacenter.Name,
			datacenter.Description,
			datacenter.Location.Name,
		)
	}
	w.Flush()

	return nil
}
