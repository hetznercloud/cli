package base

import (
	"context"
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

// AddLabelCmd allows defining commands for adding labels to resources.
type AddLabelCmd struct {
	ResourceNameSingular string
	ShortDescription     string
	NameSuggestions      func(client hcapi2.Client) func() []string
	FetchLabels          func(ctx context.Context, client hcapi2.Client, idOrName string) (map[string]string, int, error)
	SetLabels            func(ctx context.Context, client hcapi2.Client, id int, labels map[string]string) error
}

// CobraCommand creates a command that can be registered with cobra.
func (lc *AddLabelCmd) CobraCommand(
	ctx context.Context, client hcapi2.Client, tokenEnsurer state.TokenEnsurer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("add-label [FLAGS] %s LABEL", strings.ToUpper(lc.ResourceNameSingular)),
		Short:                 lc.ShortDescription,
		Args:                  cobra.ExactArgs(2),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(lc.NameSuggestions(client))),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateAddLabel, tokenEnsurer.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.Run(ctx, client, cmd, args)
		},
	}
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

// Run executes an add label command
func (lc *AddLabelCmd) Run(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	idOrName := args[0]

	labels, id, err := lc.FetchLabels(ctx, client, idOrName)
	if err != nil {
		return err
	}
	key, val := util.SplitLabelVars(args[1])

	if _, ok := labels[key]; ok && !overwrite {
		return fmt.Errorf("label %s on %s %d already exists", key, lc.ResourceNameSingular, id)
	}
	labels[key] = val

	if err := lc.SetLabels(ctx, client, id, labels); err != nil {
		return err
	}

	fmt.Printf("Label %s added to %s %d\n", key, lc.ResourceNameSingular, id)
	return nil
}

func validateAddLabel(cmd *cobra.Command, args []string) error {
	label := util.SplitLabel(args[1])
	if len(label) != 2 {
		return fmt.Errorf("invalid label: %s", args[1])
	}

	return nil
}
