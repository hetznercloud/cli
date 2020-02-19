package cli

import (
	"encoding/json"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLocationDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] LOCATION",
		Short:                 "Describe a location",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runLocationDescribe),
	}
	addOutputFlag(cmd, outputOptionJSON(), outputOptionFormat())
	return cmd
}

func runLocationDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	outputFlags := outputFlagsForCommand(cmd)

	idOrName := args[0]
	location, resp, err := cli.Client().Location.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if location == nil {
		return fmt.Errorf("location not found: %s", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return locationDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return describeFormat(location, outputFlags["format"][0])
	default:
		return locationDescribeText(cli, location)
	}
}

func locationDescribeText(cli *CLI, location *hcloud.Location) error {
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

func locationDescribeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if location, ok := data["location"]; ok {
		return describeJSON(location)
	}
	if locations, ok := data["locations"].([]interface{}); ok {
		return describeJSON(locations[0])
	}
	return describeJSON(data)
}
