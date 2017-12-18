package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newServerTypeListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "list",
		Short:            "List server types",
		TraverseChildren: true,
		RunE:             cli.wrap(runServerTypeList),
	}
	return cmd
}

func runServerTypeList(cli *CLI, cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	serverTypes, err := cli.Client().ServerType.All(ctx)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCORES\tMEMORY\tDISK\tSTORAGE TYPE")
	for _, serverType := range serverTypes {
		fmt.Fprintf(w, "%d\t%.50s\t%d\t%.1f GB\t%d GB\t%s\n",
			serverType.ID,
			serverType.Name,
			serverType.Cores,
			serverType.Memory,
			serverType.Disk,
			serverType.StorageType,
		)
	}
	w.Flush()

	return nil
}
