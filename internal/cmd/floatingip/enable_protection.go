package floatingip

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
		Use:   "enable-protection [FLAGS] FLOATINGIP PROTECTIONLEVEL [PROTECTIONLEVEL...]",
		Short: "Enable resource protection for a Floating IP",
		Args:  cobra.MinimumNArgs(2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.FloatingIPNames),
			cmpl.SuggestCandidates("delete"),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runFloatingIPEnableProtection),
	}
	return cmd
}

func runFloatingIPEnableProtection(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	floatingIP, _, err := cli.Client().FloatingIP.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if floatingIP == nil {
		return fmt.Errorf("Floating IP not found: %v", idOrName)
	}

	var unknown []string
	opts := hcloud.FloatingIPChangeProtectionOpts{}
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

	action, _, err := cli.Client().FloatingIP.ChangeProtection(cli.Context, floatingIP, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Resource protection enabled for Floating IP %d\n", floatingIP.ID)
	return nil
}
