package config

import (
	_ "embed"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

//go:embed helptext/other.txt
var nonPreferenceOptions string

//go:embed helptext/preferences.txt
var preferenceOptions string

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long: `This command allows you to manage options for the Hetzner Cloud CLI. Options can be set inside the
configuration file, through environment variables or with flags. 

The hierarchy for configuration sources is as follows (from highest to lowest priority):
1. Flags
2. Environment variables
3. Configuration file (context)
4. Configuration file (global)
5. Default values

Most options are 'preferences' - these options can be set globally and can additionally be overridden
for each context. Below is a list of all non-preference options:

` + nonPreferenceOptions +
			`
Since the above options are not preferences, they cannot be modified with 'hcloud config set' or 
'hcloud config unset'. However, you are able to retrieve them using 'hcloud config get' and 'hcloud config list'.
Following options are preferences and can be used with these commands:

` + preferenceOptions +
			`
Options will be persisted in the configuration file. To find out where your configuration file is located
on disk, run 'hcloud config get config'.
`,
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		NewSetCommand(s),
		NewGetCommand(s),
		NewListCommand(s),
		NewUnsetCommand(s),
		NewAddCommand(s),
		NewRemoveCommand(s),
	)
	return cmd
}

func gen() {

}
