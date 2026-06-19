package datacenter

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

const deprecationNotice = `The 'hcloud datacenter ...' commands are deprecated and will be removed after 1 Oct. 2026.
After this date, requests to the datacenters API endpoints will return HTTP 410 Gone.

See https://docs.hetzner.cloud/changelog#2026-06-02-datacenters-deprecated for more details.
`

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "datacenter",
		Aliases:               []string{"datacenters"},
		Short:                 "View Datacenters (deprecated)",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,

		Long: deprecationNotice,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			cmd.PrintErrln("Warning: The 'datacenter' commands are deprecated. Use the 'location' commands instead.")
		},
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
	)
	return cmd
}
