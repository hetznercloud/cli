package datacenter

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "datacenter",
		Aliases:               []string{"datacenters"},
		Short:                 "View Datacenters",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		Deprecated:            `see https://docs.hetzner.cloud/changelog#2026-06-02-datacenters-deprecated for more details.`,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
	)
	return cmd
}
