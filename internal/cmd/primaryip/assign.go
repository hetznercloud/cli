package primaryip

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var AssignCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "assign --server <server> <primary-ip>",
			Short: "Assign a Primary IP to an assignee (usually a Server)",
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.PrimaryIP().Names(true, false, nil)),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("server", "", "Name or ID of the Server")
		cmpl.RegisterFlagCompletion(cmd, "server", cmpl.SuggestCandidatesF(client.Server().Names))
		util.MarkFlagRequired(cmd, "server")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		primaryIP, _, err := s.Client().PrimaryIP().Get(cmd.Context(), idOrName)
		if err != nil {
			return err
		}
		if primaryIP == nil {
			return fmt.Errorf("Primary IP not found: %v", idOrName)
		}

		serverIDOrName, _ := cmd.Flags().GetString("server")

		server, _, err := s.Client().Server().Get(cmd.Context(), serverIDOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("Server not found: %s", serverIDOrName)
		}
		opts := hcloud.PrimaryIPAssignOpts{
			ID:           primaryIP.ID,
			AssigneeType: "server",
			AssigneeID:   server.ID,
		}

		action, _, err := s.Client().PrimaryIP().Assign(cmd.Context(), opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(cmd.Context(), cmd, action); err != nil {
			return err
		}

		cmd.Printf("Primary IP %d assigned to %s %d\n", opts.ID, cases.Title(language.English).String(opts.AssigneeType), opts.AssigneeID)
		return nil
	},
}
