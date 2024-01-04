package loadbalancer

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var DeleteServiceCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "delete-service [FLAGS] LOADBALANCER",
			Short:                 "Deletes a service from a Load Balancer",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.LoadBalancer().Names)),
			Args:                  cobra.RangeArgs(1, 2),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Int("listen-port", 0, "The listen port of the service you want to delete (required)")
		cmd.MarkFlagRequired("listen-port")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		listenPort, _ := cmd.Flags().GetInt("listen-port")
		idOrName := args[0]
		loadBalancer, _, err := s.Client().LoadBalancer().Get(s, idOrName)
		if err != nil {
			return err
		}
		if loadBalancer == nil {
			return fmt.Errorf("Load Balancer not found: %s", idOrName)
		}
		_, _, err = s.Client().LoadBalancer().DeleteService(s, loadBalancer, listenPort)
		if err != nil {
			return err
		}

		cmd.Printf("Service on port %d deleted from Load Balancer %d\n", listenPort, loadBalancer.ID)
		return nil
	},
}
