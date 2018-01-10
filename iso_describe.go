package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newISODescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] ISO",
		Short:                 "Describe a server type",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runISODescribe),
	}
	return cmd
}

func runISODescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	iso, _, err := cli.Client().ISO.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if iso == nil {
		return fmt.Errorf("iso not found: %s", idOrName)
	}

	fmt.Printf("ID:\t\t%d\n", iso.ID)
	fmt.Printf("Name:\t\t%s\n", iso.Name)
	fmt.Printf("Description:\t%s\n", iso.Description)
	fmt.Printf("Type:\t\t%s\n", iso.Type)

	return nil
}
