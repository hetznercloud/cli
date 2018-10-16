package cli

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

	return cmd
}

func runVolumeCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	opts, err := volumeOptsFromFlags(cli, cmd.Flags())
	if err != nil {
		return err
	}

	result, _, err := cli.Client().Volume.Create(cli.Context, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, result.Action); err != nil {
		return err
	}
	fmt.Printf("Volume %d created\n", result.Volume.ID)

	return nil
}

func volumeOptsFromFlags(cli *CLI, flags *pflag.FlagSet) (opts hcloud.VolumeCreateOpts, err error) {
	name, _ := flags.GetString("name")
	server, _ := flags.GetString("server")
	size, _ := flags.GetInt("size")
	location, _ := flags.GetString("location")

	opts = hcloud.VolumeCreateOpts{
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
	if server != "" {
		id, err := strconv.Atoi(server)
		if err == nil {
			opts.Server = &hcloud.Server{ID: id}
		} else {
			opts.Server = &hcloud.Server{Name: server}
		}
	}
	return
}
