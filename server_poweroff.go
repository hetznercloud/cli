package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerPoweroffCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "poweroff [flags] <id>",
		Short:            "Poweroff a server",
		Args:             cobra.ExactArgs(1),
		TraverseChildren: true,
		RunE:             cli.wrap(runServerPoweroff),
	}
	return cmd
}

func runServerPoweroff(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid server id")
	}

	ctx := context.Background()
	server := &hcloud.Server{ID: id}
	action, _, err := cli.Client().Server.Poweroff(ctx, server)
	if err != nil {
		return err
	}
	if err := <-waitAction(ctx, cli.Client(), action); err != nil {
		return err
	}
	fmt.Printf("Server %d stopped\n", id)
	return nil
}
