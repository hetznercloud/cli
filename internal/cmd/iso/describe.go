package iso

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
		Use:                   "describe [FLAGS] ISO",
		Short:                 "Describe an ISO",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ISONames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runDescribe),
	}
	output.AddFlag(cmd, output.OptionJSON(), output.OptionFormat())
	return cmd
}

func runDescribe(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := output.FlagsForCommand(cmd)

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
		return describeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(iso, outputFlags["format"][0])
	default:
		return describeText(iso)
	}
}

func describeText(iso *hcloud.ISO) error {
	fmt.Printf("ID:\t\t%d\n", iso.ID)
	fmt.Printf("Name:\t\t%s\n", iso.Name)
	fmt.Printf("Description:\t%s\n", iso.Description)
	fmt.Printf("Type:\t\t%s\n", iso.Type)
	return nil
}

func describeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if iso, ok := data["iso"]; ok {
		return util.DescribeJSON(iso)
	}
	if isos, ok := data["isos"].([]interface{}); ok {
		return util.DescribeJSON(isos[0])
	}
	return util.DescribeJSON(data)
}
