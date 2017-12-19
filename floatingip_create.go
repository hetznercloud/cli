package cli

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newFloatingIPCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "create",
		Short:            "Create a Floating IP",
		Args:             cobra.NoArgs,
		TraverseChildren: true,
		RunE:             cli.wrap(runFloatingIPCreate),
		PreRunE:          validateFloatingIPCreate,
	}
	cmd.Flags().String("type", "", "Type")
	cmd.Flags().String("description", "", "Description")
	cmd.Flags().String("home-location", "", "Home location")
	cmd.Flags().Int("server", 0, "Server to assign Floating IP to")
	return cmd
}

func validateFloatingIPCreate(cmd *cobra.Command, args []string) error {
	typ, _ := cmd.Flags().GetString("type")
	if typ == "" {
		return errors.New("type is required")
	}

	homeLocation, _ := cmd.Flags().GetString("home-location")
	server, _ := cmd.Flags().GetInt("server")
	if homeLocation == "" && server == 0 {
		return errors.New("one of --home-location or --server is required")
	}

	return nil
}

func runFloatingIPCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	typ, _ := cmd.Flags().GetString("type")
	description, _ := cmd.Flags().GetString("description")
	homeLocation, _ := cmd.Flags().GetString("home-location")
	server, _ := cmd.Flags().GetInt("server")

	opts := hcloud.FloatingIPCreateOpts{
		Type:        hcloud.FloatingIPType(typ),
		Description: &description,
	}
	if homeLocation != "" {
		opts.HomeLocation = &hcloud.Location{Name: homeLocation}
	}
	if server != 0 {
		opts.Server = &hcloud.Server{ID: server}
	}

	result, _, err := cli.Client().FloatingIP.Create(cli.Context, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Floating IP %d created\n", result.FloatingIP.ID)

	return nil
}
