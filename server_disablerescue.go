package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerDisableRescueCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "disable-rescue [flags] <id>",
		Short:            "Disable rescue for a server",
		Args:             cobra.ExactArgs(1),
		TraverseChildren: true,
		RunE:             cli.wrap(runServerDisableRescue),
	}
	return cmd
}

func runServerDisableRescue(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid server id")
	}

	var (
		ctx    = context.Background()
		server = &hcloud.Server{ID: id}
	)

	action, _, err := cli.Client().Server.DisableRescue(ctx, server)
	if err != nil {
		return err
	}
	if err := <-waitAction(ctx, cli.Client(), action); err != nil {
		return err
	}
	fmt.Printf("Rescue disabled for server %d\n", id)
	return nil
}
