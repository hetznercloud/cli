package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func newImageListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "list",
		Short:            "List images",
		TraverseChildren: true,
		RunE:             cli.wrap(runImageList),
	}
	return cmd
}

func runImageList(cli *CLI, cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	images, err := cli.Client().Image.All(ctx)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tTYPE\tNAME\tDESCRIPTION\tIMAGE SIZE\tDISK SIZE\tCREATED")
	for _, image := range images {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%.1f GB\t%.0f GB\t%s\n", image.ID, image.Type, na(image.Name),
			image.Description, image.ImageSize, image.DiskSize, humanize.Time(image.Created))
	}
	w.Flush()

	return nil
}
