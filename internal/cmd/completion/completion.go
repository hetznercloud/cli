package completion

import (
	"fmt"
	"os"

	"github.com/hetznercloud/cli/internal/state"
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

func NewCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [FLAGS] SHELL",
		Short: "Output shell completion code for the specified shell",
		Long: `To load completions:

### Bash

To load completions into the current shell execute:

    source <(hcloud completion bash)

In order to make the completions permanent, append the line above to
your .bashrc.

### Zsh

If shell completions are not already enabled for your environment need
to enable them. Add the following line to your ~/.zshrc file:

    autoload -Uz compinit; compinit

To load completions for each session execute the following commands:

    mkdir -p ~/.config/hcloud/completion/zsh
    hcloud completion zsh > ~/.config/hcloud/completion/zsh/_hcloud

Finally add the following line to your ~/.zshrc file, *before* you
call the compinit function:

    fpath+=(~/.config/hcloud/completion/zsh)

In the end your ~/.zshrc file should contain the following two lines
in the order given here.

    fpath+=(~/.config/hcloud/completion/zsh)
    #  ... anything else that needs to be done before compinit
    autoload -Uz compinit; compinit
    # ...

You will need to start a new shell for this setup to take effect.

### Fish

To load completions into the current shell execute:

    hcloud completion fish | source

In order to make the completions permanent execute once:

     hcloud completion fish > ~/.config/fish/completions/hcloud.fish
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
