package context

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
)

func newUseCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "use [FLAGS] [NAME]",
		Short:                 "Use a context",
		Args:                  cobra.MaximumNArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.Config.ContextNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runUse),
	}
	return cmd
}

func runUse(cli *state.State, cmd *cobra.Command, args []string) error {
	if os.Getenv("HCLOUD_TOKEN") != "" {
		_, _ = fmt.Fprintln(os.Stderr, "Warning: HCLOUD_TOKEN is set. The active context will have no effect.")
	}

	var name string
	if cli.Config.ActiveContext != nil {
		name = cli.Config.ActiveContext.Name
	}

	if len(args) > 0 {
		name = args[0]

	} else {
		var opts []huh.Option[string]
		for _, ctx := range cli.Config.Contexts {
			var opt huh.Option[string]
			if ctx == cli.Config.ActiveContext {
				opt = huh.NewOption(ctx.Name+" (*)", ctx.Name).Selected(true)
			} else {
				opt = huh.NewOption(ctx.Name, ctx.Name)
			}
			opts = append(opts, opt)
		}
		err := huh.NewForm(huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select context").
				Options(opts...).
				Value(&name),
		)).Run()
		if err != nil {
			return err
		}
	}

	context := cli.Config.ContextByName(name)
	if context == nil {
		return fmt.Errorf("context not found: %v", name)
	}

	cli.Config.ActiveContext = context
	return cli.WriteConfig()
}
