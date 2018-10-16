package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newVolumeAttachCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "attach [FLAGS] VOLUME",
		Short:                 "Attach a volume to a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runVolumeAttach),
	}
	cmd.Flags().String("server", "", "Server (id or name)")
	cmd.Flag("server").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_server_names"},
	}
	cmd.MarkFlagRequired("server")
	return cmd
}

func runVolumeAttach(cli *CLI, cmd *cobra.Command, args []string) error {
	volume, _, err := cli.Client().Volume.Get(cli.Context, args[0])
	if err != nil {
		return err
	}
	if volume == nil {
		return fmt.Errorf("volume not found: %s", args[0])
	}

	serverIDOrName, _ := cmd.Flags().GetString("server")
	server, _, err := cli.Client().Server.Get(cli.Context, serverIDOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", serverIDOrName)
	}

	action, _, err := cli.Client().Volume.Attach(cli.Context, volume, server)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Volume %d attached to server %s\n", volume.ID, server.Name)
	return nil
}
