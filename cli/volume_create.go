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

	cmd.Flags().String("server", "", "Server (id or name)")
	cmd.Flag("server").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_server_names"},
	}

	cmd.Flags().String("location", "", "Location (ID or name)")
	cmd.Flag("location").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_location_names"},
	}

	cmd.Flags().Int("size", 0, "Size (GB)")
	cmd.MarkFlagRequired("size")

	cmd.Flags().Bool("automount", false, "Auto mount volume after attach (Server must be provided)")
	cmd.Flags().String("format", "", "Format volume after creation (One of: xfs, ext4) (Automount must be set)")
	return cmd
}

func runVolumeCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	serverIDOrName, _ := cmd.Flags().GetString("server")
	size, _ := cmd.Flags().GetInt("size")
	location, _ := cmd.Flags().GetString("location")
	automount, _ := cmd.Flags().GetBool("automount")
	format, _ := cmd.Flags().GetString("format")

	opts := hcloud.VolumeCreateOpts{
		Name: name,
		Size: size,
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
	if automount == true {
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
