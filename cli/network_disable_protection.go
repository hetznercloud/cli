package cli

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newNetworkDisableProtectionCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "disable-protection [FLAGS] NETWORK PROTECTIONLEVEL [PROTECTIONLEVEL...]",
		Short:                 "Disable resource protection for a network",
		Args:                  cobra.MinimumNArgs(2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runNetworkDisableProtection),
	}
	return cmd
}

func runNetworkDisableProtection(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	network, _, err := cli.Client().Network.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if network == nil {
		return fmt.Errorf("network not found: %s", idOrName)
	}

	var unknown []string
	opts := hcloud.NetworkChangeProtectionOpts{}
	for _, arg := range args[1:] {
		switch strings.ToLower(arg) {
		case "delete":
			opts.Delete = hcloud.Bool(false)
		default:
			unknown = append(unknown, arg)
		}
	}
	if len(unknown) > 0 {
		return fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
	}

	action, _, err := cli.Client().Network.ChangeProtection(cli.Context, network, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Resource protection disabled for network %d\n", network.ID)
	return nil
}
