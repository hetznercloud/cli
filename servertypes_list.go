package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var serverTypeListTableOutput *tableOutput

func init() {
	serverTypeListTableOutput = newTableOutput().
		AddAllowedFields(hcloud.ServerType{}).
		AddFieldAlias("storagetype", "storage type").
		AddFieldOutputFn("memory", fieldOutputFn(func(obj interface{}) string {
			serverType := obj.(*hcloud.ServerType)
			return fmt.Sprintf("%.1f GB", serverType.Memory)
		})).
		AddFieldOutputFn("disk", fieldOutputFn(func(obj interface{}) string {
			serverType := obj.(*hcloud.ServerType)
			return fmt.Sprintf("%d GB", serverType.Disk)
		}))
}

func newServerTypeListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List server types",
		Long: listLongDescription(
			"Displays a list of server types.",
			serverTypeListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerTypeList),
	}
	return cmd
}

func runServerTypeList(cli *CLI, cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetStringArray("output")
	outOpts, err := parseOutputOpts(out)
	if err != nil {
		return err
	}

	serverTypes, err := cli.Client().ServerType.All(cli.Context)
	if err != nil {
		return err
	}

	cols := []string{"id", "name", "cores", "memory", "disk", "storage_type"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := serverTypeListTableOutput
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, serverType := range serverTypes {
		tw.Write(cols, serverType)
	}
	tw.Flush()
	return nil
}
