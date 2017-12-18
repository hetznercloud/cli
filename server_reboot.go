package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerRebootCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "reboot [flags] <id>",
		Short:            "Reboot a server",
		Args:             cobra.ExactArgs(1),
		TraverseChildren: true,
		RunE:             cli.wrap(runServerReboot),
	}
	return cmd
}

func runServerReboot(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid server id")
	}

	ctx := context.Background()
	server := &hcloud.Server{ID: id}
	action, _, err := cli.Client().Server.Reboot(ctx, server)
	if err != nil {
		return err
	}
	errCh, _ := waitAction(ctx, cli.Client(), action)
	if err := <-errCh; err != nil {
		return err
	}
	fmt.Printf("Server %d rebooted\n", id)
	return nil
}
