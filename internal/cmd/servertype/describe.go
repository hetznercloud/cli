package servertype

import (
	"encoding/json"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newDescribeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] SERVERTYPE",
		Short:                 "Describe a server type",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerTypeNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runServerTypeDescribe),
	}
	output.AddFlag(cmd, output.OptionJSON(), output.OptionFormat())
	return cmd
}

func runServerTypeDescribe(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := output.FlagsForCommand(cmd)

	idOrName := args[0]
	serverType, resp, err := cli.Client().ServerType.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if serverType == nil {
		return fmt.Errorf("server type not found: %s", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return serverTypeDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(serverType, outputFlags["format"][0])
	default:
		return serverTypeDescribeText(cli, serverType)
	}
}

func serverTypeDescribeText(cli *state.State, serverType *hcloud.ServerType) error {
	fmt.Printf("ID:\t\t%d\n", serverType.ID)
	fmt.Printf("Name:\t\t%s\n", serverType.Name)
	fmt.Printf("Description:\t%s\n", serverType.Description)
	fmt.Printf("Cores:\t\t%d\n", serverType.Cores)
	fmt.Printf("CPU Type:\t%s\n", serverType.CPUType)
	fmt.Printf("Memory:\t\t%.1f GB\n", serverType.Memory)
	fmt.Printf("Disk:\t\t%d GB\n", serverType.Disk)
	fmt.Printf("Storage Type:\t%s\n", serverType.StorageType)

	fmt.Printf("Pricings per Location:\n")
	for _, price := range serverType.Pricings {
		fmt.Printf("  - Location:\t%s:\n", price.Location.Name)
		fmt.Printf("    Hourly:\t€ %s\n", price.Hourly.Gross)
		fmt.Printf("    Monthly:\t€ %s\n", price.Monthly.Gross)
	}
	return nil
}

func serverTypeDescribeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if serverType, ok := data["server_type"]; ok {
		return util.DescribeJSON(serverType)
	}
	if serverTypes, ok := data["server_types"].([]interface{}); ok {
		return util.DescribeJSON(serverTypes[0])
	}
	return util.DescribeJSON(data)
}
