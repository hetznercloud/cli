package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerResetPasswordCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "reset-password [flags] <id>",
		Short:            "Reset password of a server",
		Args:             cobra.ExactArgs(1),
		TraverseChildren: true,
		RunE:             cli.wrap(runServerResetPassword),
	}
	return cmd
}

func runServerResetPassword(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid server id")
	}

	server := &hcloud.Server{ID: id}
	result, _, err := cli.Client().Server.ResetPassword(cli.Context, server)
	if err != nil {
		return err
	}
	errCh, _ := waitAction(cli.Context, cli.Client(), result.Action)
	if err := <-errCh; err != nil {
		return err
	}
	fmt.Printf("Password of server %d reset to: %s\n", id, result.RootPassword)
	return nil
}
