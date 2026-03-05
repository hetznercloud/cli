package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/version"
)

func NewCommand(_ state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		Short:                 "Print version information",
		Args:                  util.Validate,
		DisableFlagsInUseLine: true,
		RunE:                  runVersion,
	}
	cmd.Flags().Bool("long", false, "Print more version information (true, false)")
	return cmd
}

func runVersion(cmd *cobra.Command, _ []string) error {
	cmd.Printf("hcloud %s\n", version.Version)

	long, _ := cmd.Flags().GetBool("long")
	if long {
		tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 8, 2, ' ', 0)
		fmt.Fprintln(tw)
		fmt.Fprintf(tw, "go version:\t%s\n", runtime.Version())
		fmt.Fprintf(tw, "platform:\t%s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Fprintf(tw, "compiler:\t%s\n", runtime.Compiler)

		if info, ok := debug.ReadBuildInfo(); ok {
			fmt.Fprintf(tw, "revision:\t%s\n", getSettingsValue(info.Settings, "vcs.revision", "unknown"))
			fmt.Fprintf(tw, "revision date:\t%s\n", getSettingsValue(info.Settings, "vcs.time", "unknown"))
			fmt.Fprintf(tw, "modified:\t%s\n", getSettingsValue(info.Settings, "vcs.modified", "unknown"))
		}
		return tw.Flush()
	}
	return nil
}

func getSettingsValue(settings []debug.BuildSetting, key, def string) string {
	for _, setting := range settings {
		if setting.Key == key {
			return setting.Value
		}
	}
	return def
}
