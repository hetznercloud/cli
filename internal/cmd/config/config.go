package config

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long: `This command allows you to manage options for the Hetzner Cloud CLI. Options can be set inside the
configuration file, through environment variables or with flags. Most options are 'preferences' -
these options can be set globally and can additionally be overridden for each context. 

Below is a list of all non-preference options:

|-------------------|--------------------------------|-----------------|-------------------------|------------------|
| Option            | Description                    | Config key      | Environment variable    | Flag             |
|-------------------|--------------------------------|-----------------|-------------------------|------------------|
| config            | Config file path               |                 | HCLOUD_CONFIG           | --config         |
|-------------------|--------------------------------|-----------------|-------------------------|------------------|
| token             | API token                      | token           | HCLOUD_TOKEN            |                  |
|-------------------|--------------------------------|-----------------|-------------------------|------------------|
| context           | Currently active context       | active_context  | HCLOUD_CONTEXT          | --context        |
|-------------------|--------------------------------|-----------------|-------------------------|------------------|

Since the above options are not preferences, they cannot be modified with 'hcloud config set' or 
'hcloud config unset'. However, you are able to retrieve them using 'hcloud config get' and 'hcloud config list'.
Following options are preferences and can be used with these commands:

|-------------------|-------------------------------|------------------|-------------------------|------------------|
| Option            | Description                   | Config key       | Environment variable    | Flag             |
|-------------------|-------------------------------|------------------|-------------------------|------------------|
| endpoint          | API Endpoint to use           | endpoint         | HCLOUD_ENDPOINT         | --endpoint       |
|-------------------|-------------------------------|------------------|-------------------------|------------------|
| debug             | Enable debug output           | debug            | HCLOUD_DEBUG            | --debug          |
|-------------------|-------------------------------|------------------|-------------------------|------------------|
| debug-file        | File to write debug output to | debug_file       | HCLOUD_DEBUG_FILE       | --debug-file     |
|-------------------|-------------------------------|------------------|-------------------------|------------------|
| poll-interval     | Time between requests         | endpoint         | HCLOUD_POLL_INTERVAL    | --poll-interval  |
|                   | when polling                  |                  |                         |                  |
|-------------------|-------------------------------|------------------|-------------------------|------------------|
| quiet             | If true, only error messages  | quiet            | HCLOUD_QUIET            | --quiet          |
|                   | are printed                   |                  |                         |                  |
|-------------------|-------------------------------|------------------|-------------------------|------------------|
| default-ssh-keys  | Default SSH keys for server   | default_ssh_keys | HCLOUD_DEFAULT_SSH_KEYS |                  |
|                   | creation                      |                  |                         |                  |
|-------------------|-------------------------------|------------------|-------------------------|------------------|
| ssh-path          | Path to the SSH binary (used  | ssh_path         | HCLOUD_SSH_PATH         | --ssh-path       |
|                   | for 'hcloud server ssh')      |                  |                         |                  |
|-------------------|-------------------------------|------------------|-------------------------|------------------|

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
