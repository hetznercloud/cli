package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerPoweronCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "poweron [flags] <id>",
		Short:            "Poweron a server",
		Args:             cobra.ExactArgs(1),
		TraverseChildren: true,
		RunE:             cli.wrap(runServerPoweron),
	}
	return cmd
}

func runServerPoweron(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid server id")
	}

	ctx := context.Background()
	server := &hcloud.Server{ID: id}
	action, _, err := cli.Client().Server.Poweron(ctx, server)
	if err != nil {
		return err
	}
	errCh, _ := waitAction(ctx, cli.Client(), action)
	if err := <-errCh; err != nil {
		return err
	}
	fmt.Printf("Server %d started\n", id)
	return nil
}
