package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerShutdownCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "shutdown [flags] <id>",
		Short:            "Shutdown a server",
		Args:             cobra.ExactArgs(1),
		TraverseChildren: true,
		RunE:             cli.wrap(runServerShutdown),
	}
	return cmd
}

func runServerShutdown(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid server id")
	}

	server := &hcloud.Server{ID: id}
	action, _, err := cli.Client().Server.Shutdown(cli.Context, server)
	if err != nil {
		return err
	}
	errCh, _ := waitAction(cli.Context, cli.Client(), action)
	if err := <-errCh; err != nil {
		return err
	}
	fmt.Printf("Server %d shut down\n", id)
	return nil
}
