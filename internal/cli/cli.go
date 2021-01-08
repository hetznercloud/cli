package cli

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

type CLI struct {
	State       *state.State
	RootCommand *cobra.Command
}

func NewCLI() *CLI {
	cli := &CLI{
		State: state.New(),
	}

	cli.RootCommand = NewRootCommand(cli.State)
	return cli
}
