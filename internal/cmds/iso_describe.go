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

func newISODescribeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] ISO",
		Short:                 "Describe an ISO",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ISONames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runISODescribe),
	}
	addOutputFlag(cmd, outputOptionJSON(), outputOptionFormat())
	return cmd
}

func runISODescribe(cli *state.State, cmd *cobra.Command, args []string) error {
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
		return util.DescribeFormat(iso, outputFlags["format"][0])
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
		return util.DescribeJSON(iso)
	}
	if isos, ok := data["isos"].([]interface{}); ok {
		return util.DescribeJSON(isos[0])
	}
	return util.DescribeJSON(data)
}
