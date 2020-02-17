package cli

import (
	"encoding/json"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newISODescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] ISO",
		Short:                 "Describe an ISO",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runISODescribe),
	}
	addOutputFlag(cmd, outputOptionJSON(), outputOptionFormat())
	return cmd
}

func runISODescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	outputFlags := outputFlagsForCommand(cmd)

	idOrName := args[0]
	iso, resp, err := cli.Client().ISO.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if iso == nil {
		return fmt.Errorf("iso not found: %s", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return isoDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return describeFormat(iso, outputFlags["format"][0])
	default:
		return isoDescribeText(iso)
	}
}

func isoDescribeText(iso *hcloud.ISO) error {
	fmt.Printf("ID:\t\t%d\n", iso.ID)
	fmt.Printf("Name:\t\t%s\n", iso.Name)
	fmt.Printf("Description:\t%s\n", iso.Description)
	fmt.Printf("Type:\t\t%s\n", iso.Type)
	return nil
}

func isoDescribeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if iso, ok := data["iso"]; ok {
		return describeJSON(iso, true)
	}
	if isos, ok := data["isos"].([]interface{}); ok {
		return describeJSON(isos[0], true)
	}
	return describeJSON(data, true)
}
