package version

import (
	"runtime"
	"runtime/debug"

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
		m := versionMetadata()
		cmd.Printf(`revision:   %s
build date: %s
go version: %s
platform:   %s
`,
			m["revision"],
			m["build date"],
			m["go version"],
			m["platform"],
		)
	}
	return nil
}

func versionMetadata() map[string]string {
	m := map[string]string{
		"revision":   "unknown",
		"build date": "unknown",
		"go version": runtime.Version(),
		"platform":   runtime.GOOS + "/" + runtime.GOARCH,
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return m
	}

	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			m["revision"] = setting.Value
		}
		if setting.Key == "vcs.time" {
			m["date"] = setting.Value
		}
	}

	return m
}
