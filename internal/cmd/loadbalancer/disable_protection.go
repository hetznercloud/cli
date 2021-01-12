package loadbalancer

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newDisableProtectionCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-protection [FLAGS] LOADBALANCER PROTECTIONLEVEL [PROTECTIONLEVEL...]",
		Short: "Disable resource protection for a Load Balancer",
		Args:  cobra.MinimumNArgs(2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.LoadBalancerNames),
			cmpl.SuggestCandidates("delete"),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runDisableProtection),
	}
	return cmd
}

func runDisableProtection(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}

	var unknown []string
	opts := hcloud.LoadBalancerChangeProtectionOpts{}
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

	action, _, err := cli.Client().LoadBalancer.ChangeProtection(cli.Context, loadBalancer, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Resource protection disabled for Load Balancer %d\n", loadBalancer.ID)
	return nil
}
