package firewall

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newUpdateCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "update [FLAGS] FIREWALL",
		Short:                 "Update a Firewall",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.FirewallNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runFirewallUpdate),
	}

	cmd.Flags().String("name", "", "Firewall name")

	return cmd
}

func runFirewallUpdate(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	firewall, _, err := cli.Client().Firewall.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if firewall == nil {
		return fmt.Errorf("Firewall not found: %s", idOrName)
	}

	name, _ := cmd.Flags().GetString("name")
	opts := hcloud.FirewallUpdateOpts{
		Name: name,
	}
	_, _, err = cli.Client().Firewall.Update(cli.Context, firewall, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Firewall %d updated\n", firewall.ID)
	return nil
}
