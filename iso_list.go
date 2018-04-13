package cli

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var isoListTableOutput *tableOutput

func init() {
	isoListTableOutput = newTableOutput().
		AddAllowedFields(hcloud.ISO{})
}

func newISOListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List ISOs",
		Long: fmt.Sprintf(`Displays a list of ISOs.

%s

Columns:
 - %s`, OutputDescription, strings.Join(isoListTableOutput.Columns(), "\n - ")),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runISOList),
	}
	return cmd
}

func runISOList(cli *CLI, cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetStringArray("output")
	outOpts, err := parseOutputOpts(out)
	if err != nil {
		return err
	}

	isos, err := cli.Client().ISO.All(cli.Context)
	if err != nil {
		return err
	}

	cols := []string{"id", "name", "description", "type"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := isoListTableOutput
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, iso := range isos {
		tw.Write(cols, iso)
	}
	tw.Flush()

	return nil
}
