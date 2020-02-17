package cli

import (
	"encoding/json"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerTypeDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] SERVERTYPE",
		Short:                 "Describe a server type",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerTypeDescribe),
	}
	addOutputFlag(cmd, outputOptionJSON(), outputOptionFormat())
	return cmd
}

func runServerTypeDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	outputFlags := outputFlagsForCommand(cmd)

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
		return describeFormat(serverType, outputFlags["format"][0])
	default:
		return serverTypeDescribeText(cli, serverType)
	}
}

func serverTypeDescribeText(cli *CLI, serverType *hcloud.ServerType) error {
	fmt.Printf("ID:\t\t%d\n", serverType.ID)
	fmt.Printf("Name:\t\t%s\n", serverType.Name)
	fmt.Printf("Description:\t%s\n", serverType.Description)
	fmt.Printf("Cores:\t\t%d\n", serverType.Cores)
	fmt.Printf("Memory:\t\t%.1f GB\n", serverType.Memory)
	fmt.Printf("Disk:\t\t%d GB\n", serverType.Disk)
	fmt.Printf("Storage Type:\t%s\n", serverType.StorageType)
	return nil
}

func serverTypeDescribeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if serverType, ok := data["server_type"]; ok {
		return describeJSON(serverType, true)
	}
	if serverTypes, ok := data["server_types"].([]interface{}); ok {
		return describeJSON(serverTypes[0], true)
	}
	return describeJSON(data, true)
}
