package cli

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newFloatingIPCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create FLAGS",
		Short:                 "Create a Floating IP",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(validateFloatingIPCreate, cli.ensureToken),
		RunE:                  cli.wrap(runFloatingIPCreate),
	}
	cmd.Flags().String("type", "", "Type (ipv4 or ipv6)")
	cmd.Flag("type").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_floatingip_types"},
	}
	cmd.MarkFlagRequired("type")

	cmd.Flags().String("description", "", "Description")
	cmd.Flags().String("name", "", "Name")
	cmd.Flags().String("home-location", "", "Home location")
	cmd.Flag("home-location").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_location_names"},
	}

	cmd.Flags().String("server", "", "Server to assign Floating IP to")
	cmd.Flag("server").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_server_names"},
	}

	cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

	return cmd
}

func validateFloatingIPCreate(cmd *cobra.Command, args []string) error {
	typ, _ := cmd.Flags().GetString("type")
	if typ == "" {
		return errors.New("type is required")
	}

	homeLocation, _ := cmd.Flags().GetString("home-location")
	server, _ := cmd.Flags().GetString("server")
	if homeLocation == "" && server == "" {
		return errors.New("one of --home-location or --server is required")
	}

	return nil
}

func runFloatingIPCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	typ, _ := cmd.Flags().GetString("type")
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	homeLocation, _ := cmd.Flags().GetString("home-location")
	serverNameOrID, _ := cmd.Flags().GetString("server")
	labels, _ := cmd.Flags().GetStringToString("label")

	opts := hcloud.FloatingIPCreateOpts{
		Type:        hcloud.FloatingIPType(typ),
		Description: &description,
		Labels:      labels,
	}
	if name != "" {
		opts.Name = &name
	}
	if homeLocation != "" {
		opts.HomeLocation = &hcloud.Location{Name: homeLocation}
	}
	if serverNameOrID != "" {
		server, _, err := cli.Client().Server.Get(cli.Context, serverNameOrID)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", serverNameOrID)
		}
		opts.Server = server
	}

	result, _, err := cli.Client().FloatingIP.Create(cli.Context, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Floating IP %d created\n", result.FloatingIP.ID)

	return nil
}
