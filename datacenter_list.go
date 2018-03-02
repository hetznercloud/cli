package cli

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newDatacenterListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list [FLAGS]",
		Short:                 "List datacenters",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runDatacenterList),
	}
	return cmd
}

func runDatacenterList(cli *CLI, cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetStringArray("output")
	outOpts, err := parseOutputOpts(out)
	if err != nil {
		return err
	}

	datacenters, err := cli.Client().Datacenter.All(cli.Context)
	if err != nil {
		return err
	}

	cols := []string{"id", "name", "description", "location"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := newTableOutput().
		AddAllowedFields(hcloud.Datacenter{}).
		AddFieldOutputFn("location", fieldOutputFn(func(obj interface{}) string {
			datacenter := obj.(*hcloud.Datacenter)
			return datacenter.Location.Name
		}))

	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, datacenter := range datacenters {
		tw.Write(cols, datacenter)
	}
	tw.Flush()
	return nil
}
