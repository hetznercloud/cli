package cli

import "github.com/spf13/cobra"

const (
	bashCompletionFunc = `
  __hcloud_server_ids() {
    local ctl_output out
    if ctl_output=$(hcloud server list --no-header 2>/dev/null); then
        COMPREPLY=($(echo "${ctl_output}" | awk '{print $1}'))
    fi
  }

  __custom_func() {
    case ${last_command} in
      hcloud_server_delete | hcloud_server_describe )
        __hcloud_server_ids
        return
        ;;
      *)
        ;;
    esac
  }
  `
)

func NewRootCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "hcloud",
		Short:                  "Hetzner Cloud CLI",
		Long:                   "A command-line interface for Hetzner Cloud",
		RunE:                   cli.wrap(runRoot),
		TraverseChildren:       true,
		SilenceUsage:           true,
		SilenceErrors:          true,
		BashCompletionFunction: bashCompletionFunc,
	}
	cmd.Flags().BoolVar(&cli.JSON, "json", false, "Output JSON API response")
	cmd.AddCommand(
		newConfigureCommand(cli),
		newFloatingIPCommand(cli),
		newServerCommand(cli),
		newSSHKeyCommand(cli),
		newVersionCommand(cli),
		newCompletionCommand(cli),
	)
	return cmd
}

func runRoot(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
