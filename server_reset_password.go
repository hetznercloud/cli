package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newServerResetPasswordCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "reset-password [FLAGS] SERVER",
		Short:                 "Reset password of a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServerResetPassword),
	}
	return cmd
}

func runServerResetPassword(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	result, _, err := cli.Client().Server.ResetPassword(cli.Context, server)
	if err != nil {
		return err
	}
	errCh, _ := waitAction(cli.Context, cli.Client(), result.Action)
	if err := <-errCh; err != nil {
		return err
	}
	fmt.Printf("Password of server %s reset to: %s\n", idOrName, result.RootPassword)
	return nil
}
