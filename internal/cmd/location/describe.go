package location

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

type describer struct {
	client hcapi2.Client
}

func newDescribeCommand(
	ctx context.Context,
	client hcapi2.Client,
	tokenEnsurer state.TokenEnsurer,
	actionWaiter state.ActionWaiter,
) *cobra.Command {
	d := describer{
		client: client,
	}

	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] LOCATION",
		Short:                 "Describe a location",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Location().Names)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               tokenEnsurer.EnsureToken,
		RunE:                  state.WrapCtx(ctx, d.describe),
	}
	output.AddFlag(cmd, output.OptionJSON(), output.OptionFormat())
	return cmd
}

func (d *describer) describe(ctx context.Context, cmd *cobra.Command, args []string) error {
	outputFlags := output.FlagsForCommand(cmd)

	idOrName := args[0]
	location, resp, err := d.client.Location().Get(ctx, idOrName)
	if err != nil {
		return err
	}
	if location == nil {
		return fmt.Errorf("location not found: %s", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return describeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(location, outputFlags["format"][0])
	default:
		return describeText(location)
	}
}

func describeText(location *hcloud.Location) error {
	fmt.Printf("ID:\t\t%d\n", location.ID)
	fmt.Printf("Name:\t\t%s\n", location.Name)
	fmt.Printf("Description:\t%s\n", location.Description)
	fmt.Printf("Network Zone:\t%s\n", location.NetworkZone)
	fmt.Printf("Country:\t%s\n", location.Country)
	fmt.Printf("City:\t\t%s\n", location.City)
	fmt.Printf("Latitude:\t%f\n", location.Latitude)
	fmt.Printf("Longitude:\t%f\n", location.Longitude)
	return nil
}

func describeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if location, ok := data["location"]; ok {
		return util.DescribeJSON(location)
	}
	if locations, ok := data["locations"].([]interface{}); ok {
		return util.DescribeJSON(locations[0])
	}
	return util.DescribeJSON(data)
}
