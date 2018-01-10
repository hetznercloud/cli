package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newLocationListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list [FLAGS]",
		Short:                 "List locations",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runLocationList),
	}
	return cmd
}

func runLocationList(cli *CLI, cmd *cobra.Command, args []string) error {
	locations, err := cli.Client().Location.All(cli.Context)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION\tCOUNTRY\tCITY")
	for _, location := range locations {
		fmt.Fprintf(w, "%d\t%.50s\t%.50s\t%s\t%s\n",
			location.ID,
			location.Name,
			location.Description,
			location.Country,
			location.City,
		)
	}
	w.Flush()

	return nil
}
