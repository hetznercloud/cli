package cli

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newVolumeCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create FLAGS",
		Short:                 "Create a volume",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runVolumeCreate),
	}
	cmd.Flags().String("name", "", "Volume name")
	cmd.MarkFlagRequired("name")

	cmd.Flags().String("server", "", "Server (ID or name)")
	cmd.Flag("server").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_server_names"},
	}

	cmd.Flags().String("location", "", "Location (ID or name)")
	cmd.Flag("location").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_location_names"},
	}

	cmd.Flags().Int("size", 0, "Size (GB)")
	cmd.MarkFlagRequired("size")

	cmd.Flags().Bool("automount", false, "Automount volume after attach (server must be provided)")
	cmd.Flags().String("format", "", "Format volume after creation (automount must be enabled)")

	cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

	return cmd
}

func runVolumeCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	serverIDOrName, _ := cmd.Flags().GetString("server")
	size, _ := cmd.Flags().GetInt("size")
	location, _ := cmd.Flags().GetString("location")
	automount, _ := cmd.Flags().GetBool("automount")
	format, _ := cmd.Flags().GetString("format")
	labels, _ := cmd.Flags().GetStringToString("label")

	opts := hcloud.VolumeCreateOpts{
		Name:   name,
		Size:   size,
		Labels: labels,
	}

	if location != "" {
		id, err := strconv.Atoi(location)
		if err == nil {
			opts.Location = &hcloud.Location{ID: id}
		} else {
			opts.Location = &hcloud.Location{Name: location}
		}
	}
	if serverIDOrName != "" {
		server, _, err := cli.Client().Server.Get(cli.Context, serverIDOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", serverIDOrName)
		}
		opts.Server = server
	}
	if automount {
		opts.Automount = &automount
		if format != "" {
			opts.Format = &format
		}
	}

	result, _, err := cli.Client().Volume.Create(cli.Context, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, result.Action); err != nil {
		return err
	}
	if err := cli.WaitForActions(cli.Context, result.NextActions); err != nil {
		return err
	}
	fmt.Printf("Volume %d created\n", result.Volume.ID)

	return nil
}
