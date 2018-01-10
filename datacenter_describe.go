package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newDatacenterDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] DATACENTER",
		Short:                 "Describe a datacenter",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runDatacenterDescribe),
	}
	return cmd
}

func runDatacenterDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	datacenter, _, err := cli.Client().Datacenter.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if datacenter == nil {
		return fmt.Errorf("datacenter not found: %s", idOrName)
	}

	fmt.Printf("ID:\t\t%d\n", datacenter.ID)
	fmt.Printf("Name:\t\t%s\n", datacenter.Name)
	fmt.Printf("Description:\t%s\n", datacenter.Description)
	fmt.Printf("Location:\n")
	fmt.Printf("  Name:\t\t%s\n", datacenter.Location.Name)
	fmt.Printf("  Description:\t%s\n", datacenter.Location.Description)
	fmt.Printf("  Country:\t%s\n", datacenter.Location.Country)
	fmt.Printf("  City:\t\t%s\n", datacenter.Location.City)
	fmt.Printf("  Latitude:\t%f\n", datacenter.Location.Latitude)
	fmt.Printf("  Longitude:\t%f\n", datacenter.Location.Longitude)
	fmt.Printf("Server Types:\n")

	serverTypesMap := map[int]*hcloud.ServerType{}
	for _, t := range datacenter.ServerTypes.Available {
		serverTypesMap[t.ID] = t
	}
	for _, t := range datacenter.ServerTypes.Supported {
		serverTypesMap[t.ID] = t
	}
	for id := range serverTypesMap {
		var err error
		serverTypesMap[id], _, err = cli.client.ServerType.GetByID(cli.Context, id)
		if err != nil {
			return fmt.Errorf("error fetching server type: %v", err)
		}
	}

	printServerTypes := func(list []*hcloud.ServerType, dataMap map[int]*hcloud.ServerType) {
		for _, t := range list {
			st := dataMap[t.ID]
			fmt.Printf("  - ID:\t\t %d\n", st.ID)
			fmt.Printf("    Name:\t %s\n", st.Name)
			fmt.Printf("    Description: %s\n", st.Description)
		}
	}

	fmt.Printf("  Available:\n")
	if len(datacenter.ServerTypes.Available) > 0 {
		printServerTypes(datacenter.ServerTypes.Available, serverTypesMap)
	} else {
		fmt.Printf("    No available server types\n")
	}
	fmt.Printf("  Supported:\n")
	if len(datacenter.ServerTypes.Supported) > 0 {
		printServerTypes(datacenter.ServerTypes.Available, serverTypesMap)
	} else {
		fmt.Printf("    No supported server types\n")
	}

	return nil
}
