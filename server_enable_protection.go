package cli

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerEnableProtectionCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "enable-protection [FLAGS] SERVER [PROTECTION LEVELS]",
		Short:                 "Enable resource protection for a server",
		Args:                  cobra.MinimumNArgs(2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerEnableProtection),
	}
	return cmd
}

func runServerEnableProtection(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	opts := hcloud.ServerChangeProtectionOpts{}
	for _, arg := range args[1:] {
		switch strings.ToLower(arg) {
		case "delete":
			opts.Delete = hcloud.Bool(true)
		case "rebuild":
			opts.Rebuild = hcloud.Bool(true)
		}
	}

	action, _, err := cli.Client().Server.ChangeProtection(cli.Context, server, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Resource protection enabled for server %s\n", server.Name)
	return nil
}
