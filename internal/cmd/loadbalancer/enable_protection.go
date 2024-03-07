package loadbalancer

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func getChangeProtectionOpts(enable bool, flags []string) (hcloud.LoadBalancerChangeProtectionOpts, error) {

	opts := hcloud.LoadBalancerChangeProtectionOpts{}

	var unknown []string
	for _, arg := range flags {
		switch strings.ToLower(arg) {
		case "delete":
			opts.Delete = hcloud.Ptr(enable)
		default:
			unknown = append(unknown, arg)
		}
	}
	if len(unknown) > 0 {
		return opts, fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
	}

	return opts, nil
}

func changeProtection(s state.State, cmd *cobra.Command,
	loadBalancer *hcloud.LoadBalancer, enable bool, opts hcloud.LoadBalancerChangeProtectionOpts) error {

	if opts.Delete == nil {
		return nil
	}

	action, _, err := s.Client().LoadBalancer().ChangeProtection(s, loadBalancer, opts)
	if err != nil {
		return err
	}

	if err := s.ActionProgress(cmd, s, action); err != nil {
		return err
	}

	if enable {
		cmd.Printf("Resource protection enabled for Load Balancer %d\n", loadBalancer.ID)
	} else {
		cmd.Printf("Resource protection disabled for Load Balancer %d\n", loadBalancer.ID)
	}
	return nil
}

var EnableProtectionCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "enable-protection <load-balancer> <protection-level>...",
			Short: "Enable resource protection for a Load Balancer",
			Args:  util.Validate,
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.LoadBalancer().Names),
				cmpl.SuggestCandidates("delete"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		loadBalancer, _, err := s.Client().LoadBalancer().Get(s, idOrName)
		if err != nil {
			return err
		}
		if loadBalancer == nil {
			return fmt.Errorf("Load Balancer not found: %s", idOrName)
		}

		opts, err := getChangeProtectionOpts(true, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(s, cmd, loadBalancer, true, opts)
	},
}
