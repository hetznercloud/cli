package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newServerTypeDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] SERVERTYPE",
		Short:                 "Describe a server type",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServerTypeDescribe),
	}
	return cmd
}

func runServerTypeDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid server type id")
	}

	serverType, _, err := cli.Client().ServerType.Get(cli.Context, id)
	if err != nil {
		return err
	}
	if serverType == nil {
		return fmt.Errorf("server type not found: %d", id)
	}

	fmt.Printf("ID:\t\t%d\n", serverType.ID)
	fmt.Printf("Name:\t\t%s\n", serverType.Name)
	fmt.Printf("Description:\t%s\n", serverType.Description)
	fmt.Printf("Cores:\t\t%d\n", serverType.Cores)
	fmt.Printf("Memory:\t\t%.1f GB\n", serverType.Memory)
	fmt.Printf("Disk:\t\t%d GB\n", serverType.Disk)
	fmt.Printf("Storage Type:\t%s\n", serverType.StorageType)

	return nil
}
