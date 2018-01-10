package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newISOListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list [FLAGS]",
		Short:                 "List isos",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runISOList),
	}
	return cmd
}

func runISOList(cli *CLI, cmd *cobra.Command, args []string) error {
	isos, err := cli.Client().ISO.All(cli.Context)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION\tTYPE")
	for _, iso := range isos {
		fmt.Fprintf(w, "%d\t%.50s\t%s\t%s\n",
			iso.ID,
			iso.Name,
			iso.Description,
			iso.Type,
		)
	}
	w.Flush()

	return nil
}
