package primaryip

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeDNSCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:              "set-rdns [FLAGS] PRIMARYIP",
			Short:            "Change the reverse DNS from a Primary IP",
			Args:             cobra.ExactArgs(1),
			TraverseChildren: true,
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.PrimaryIP().Names),
			),
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("hostname", "", "Hostname to set as a reverse DNS PTR entry (required)")
		cmd.MarkFlagRequired("hostname")
		cmd.Flags().String("ip", "", "IP address for which the reverse DNS entry should be set (required)")
		cmd.MarkFlagRequired("ip")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		primaryIP, _, err := s.PrimaryIP().Get(s, idOrName)
		if err != nil {
			return err
		}
		if primaryIP == nil {
			return fmt.Errorf("Primary IP not found: %v", idOrName)
		}

		DNSPtr, _ := cmd.Flags().GetString("hostname")
		ip, _ := cmd.Flags().GetString("ip")
		opts := hcloud.PrimaryIPChangeDNSPtrOpts{
			ID:     primaryIP.ID,
			DNSPtr: DNSPtr,
			IP:     ip,
		}

		action, _, err := s.PrimaryIP().ChangeDNSPtr(s, opts)
		if err != nil {
			return err
		}

		if err := s.ActionProgress(cmd, s, action); err != nil {
			return err
		}

		cmd.Printf("Primary IP %d DNS pointer: %s associated to %s\n", opts.ID, opts.DNSPtr, opts.IP)
		return nil
	},
}
