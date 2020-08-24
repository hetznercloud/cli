package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// const (
// 	completionShortDescription = "Output shell completion code for the specified shell (bash or zsh)"
// 	completionLongDescription  = completionShortDescription + `

// Note: this requires the bash-completion framework, which is not installed by default on Mac. This can be installed by using homebrew:

// 	$ brew install bash-completion

// Once installed, bash completion must be evaluated. This can be done by adding the following line to the .bash profile:

// 	$ source $(brew --prefix)/etc/bash_completion

// Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2

// Examples:
// 	# Load the hcloud completion code for bash into the current shell
// 	source <(hcloud completion bash)

// 	# Load the hcloud completion code for zsh into the current shell
// 	source <(hcloud completion zsh)`
// )

func newCompletionCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [FLAGS] SHELL",
		Short: "Output shell completion code for the specified shell",
		Long: `To load completions:

Bash:

$ source <(hcloud completion bash)

# To load completions for each session, execute once:
Linux:
  $ hcloud completion bash > /etc/bash_completion.d/hcloud
MacOS:
  $ hcloud completion bash > /usr/local/etc/bash_completion.d/hcloud

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ hcloud completion zsh > "${fpath[1]}/_hcloud"

# You will need to start a new shell for this setup to take effect.

Fish:

$ hcloud completion fish | source

# To load completions for each session, execute once:
$ hcloud completion fish > ~/.config/fish/completions/hcloud.fish
`,
		Args:                  cobra.ExactArgs(1),
		ValidArgs:             []string{"bash", "fish", "zsh"},
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			switch args[0] {
			case "bash":
				err = cmd.Root().GenBashCompletion(os.Stdout)
			case "fish":
				err = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "zsh":
				err = cmd.Root().GenZshCompletion(os.Stdout)
			default:
				err = fmt.Errorf("Unsupported shell: %s", args[0])
			}
			return err
		},
	}
	return cmd
}
