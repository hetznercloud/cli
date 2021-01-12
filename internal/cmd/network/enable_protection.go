package network

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newEnableProtectionCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable-protection [FLAGS] NETWORK PROTECTIONLEVEL [PROTECTIONLEVEL...]",
		Short: "Enable resource protection for a network",
		Args:  cobra.MinimumNArgs(2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.NetworkNames),
			cmpl.SuggestCandidates("delete"),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runEnableProtection),
	}
	return cmd
}

func runEnableProtection(cli *state.State, cmd *cobra.Command, args []string) error {
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
			opts.Delete = hcloud.Bool(true)
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

	fmt.Printf("Resource protection enabled for network %d\n", network.ID)
	return nil
}
