package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newServerDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "delete [flags] <id>",
		Short:            "Delete a server",
		Args:             cobra.ExactArgs(1),
		TraverseChildren: true,
		RunE:             cli.wrap(runServerDelete),
	}
	return cmd
}

func runServerDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid server id")
	}

	_, err = cli.Client().Server.Delete(cli.Context, id)
	if err != nil {
		return err
	}

	fmt.Printf("Server %d deleted\n", id)
	return nil
}
