package cmds

import (
	"encoding/json"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newDatacenterDescribeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] DATACENTER",
		Short:                 "Describe a datacenter",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.DataCenterNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runDatacenterDescribe),
	}
	addOutputFlag(cmd, outputOptionJSON(), outputOptionFormat())
	return cmd
}

func runDatacenterDescribe(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := outputFlagsForCommand(cmd)

	idOrName := args[0]
	datacenter, resp, err := cli.Client().Datacenter.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if datacenter == nil {
		return fmt.Errorf("datacenter not found: %s", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return datacenterDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(datacenter, outputFlags["format"][0])
	default:
		return datacenterDescribeText(cli, datacenter)
	}
}

func datacenterDescribeText(cli *state.State, datacenter *hcloud.Datacenter) error {
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
		serverTypesMap[id], _, err = cli.Client().ServerType.GetByID(cli.Context, id)
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
		printServerTypes(datacenter.ServerTypes.Supported, serverTypesMap)
	} else {
		fmt.Printf("    No supported server types\n")
	}

	return nil
}

func datacenterDescribeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if datacenter, ok := data["datacenter"]; ok {
		return util.DescribeJSON(datacenter)
	}
	if datacenters, ok := data["datacenters"].([]interface{}); ok {
		return util.DescribeJSON(datacenters[0])
	}
	return util.DescribeJSON(data)
}
