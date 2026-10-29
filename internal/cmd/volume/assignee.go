package volume

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const assigneeDeprecationNotice = `The --server flag is deprecated. Use --assignee and --assignee-type instead.
See https://docs.hetzner.cloud/changelog#2026-10-27-volumes-allow-generic-resource-targets.`

func addAssigneeFlags(cmd *cobra.Command, client hcapi2.Client) {
	cmd.Flags().String("assignee", "", "Assignee (ID, or name if the assignee type is server)")
	_ = cmd.RegisterFlagCompletionFunc("assignee", cmpl.SuggestCandidatesF(client.Server().Names))

	cmd.Flags().String("assignee-type", "", "Assignee type (server)")
	_ = cmd.RegisterFlagCompletionFunc("assignee-type", cmpl.SuggestCandidates("server"))

	cmd.Flags().String("server", "", "Server (ID or name) (deprecated)")
	_ = cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(client.Server().Names))

	cmd.MarkFlagsRequiredTogether("assignee", "assignee-type")
	cmd.MarkFlagsMutuallyExclusive("server", "assignee")
	cmd.MarkFlagsMutuallyExclusive("server", "assignee-type")
}

// getAssignee parses the command flags and returns a Volume assignee and a display name
func getAssignee(s state.State, cmd *cobra.Command) (*hcloud.VolumeAssignee, string, error) {
	assigneeIDOrName, _ := cmd.Flags().GetString("assignee")
	assigneeType, _ := cmd.Flags().GetString("assignee-type")

	if cmd.Flags().Changed("server") {
		cmd.PrintErrln("Warning: The --server flag is deprecated. Use --assignee and --assignee-type instead.")
		assigneeIDOrName, _ = cmd.Flags().GetString("server")
		assigneeType = "server"
	}

	if assigneeIDOrName == "" {
		return nil, "", nil
	}

	switch assigneeType {
	case "server":
		server, _, err := s.Client().Server().Get(s, assigneeIDOrName)
		if err != nil {
			return nil, "", err
		}
		if server == nil {
			return nil, "", fmt.Errorf("Server not found: %s", assigneeIDOrName)
		}
		return &hcloud.VolumeAssignee{Type: assigneeType, ID: server.ID}, fmt.Sprintf("Server %s", server.Name), nil
	default:
		id, err := strconv.ParseInt(assigneeIDOrName, 10, 64)
		if err != nil {
			return nil, "", fmt.Errorf("invalid assignee ID: %s", assigneeIDOrName)
		}
		return &hcloud.VolumeAssignee{Type: assigneeType, ID: id}, util.ResourceDisplayName(assigneeType, id), nil
	}
}
