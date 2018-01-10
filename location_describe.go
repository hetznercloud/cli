package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newLocationDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] LOCATION",
		Short:                 "Describe a location",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runLocationDescribe),
	}
	return cmd
}

func runLocationDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	location, _, err := cli.Client().Location.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if location == nil {
		return fmt.Errorf("location not found: %s", idOrName)
	}

	fmt.Printf("ID:\t\t%d\n", location.ID)
	fmt.Printf("Name:\t\t%s\n", location.Name)
	fmt.Printf("Description:\t%s\n", location.Description)
	fmt.Printf("Country:\t%s\n", location.Country)
	fmt.Printf("City:\t\t%s\n", location.City)
	fmt.Printf("Latitude:\t%f\n", location.Latitude)
	fmt.Printf("Longitude:\t%f\n", location.Longitude)

	return nil
}
