package base

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

// LabelCmds allows defining commands for adding labels to resources.
type LabelCmds struct {
	ResourceNameSingular   string
	ShortDescriptionAdd    string
	ShortDescriptionRemove string
	NameSuggestions        func(client hcapi2.Client) func() []string
	LabelKeySuggestions    func(client hcapi2.Client) func(idOrName string) []string
	FetchLabels            func(ctx context.Context, client hcapi2.Client, idOrName string) (map[string]string, int, error)
	SetLabels              func(ctx context.Context, client hcapi2.Client, id int, labels map[string]string) error
}

// AddCobraCommand creates a command that can be registered with cobra.
func (lc *LabelCmds) AddCobraCommand(
	ctx context.Context, client hcapi2.Client, tokenEnsurer state.TokenEnsurer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("add-label [FLAGS] %s LABEL", strings.ToUpper(lc.ResourceNameSingular)),
		Short:                 lc.ShortDescriptionAdd,
		Args:                  cobra.ExactArgs(2),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(lc.NameSuggestions(client))),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateAddLabel, tokenEnsurer.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.RunAdd(ctx, client, cmd, args)
		},
	}
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

// RunAdd executes an add label command
func (lc *LabelCmds) RunAdd(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	idOrName := args[0]

	labels, id, err := lc.FetchLabels(ctx, client, idOrName)
	if err != nil {
		return err
	}

	if labels == nil {
		labels = map[string]string{}
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

// RemoveCobraCommand creates a command that can be registered with cobra.
func (lc *LabelCmds) RemoveCobraCommand(
	ctx context.Context, client hcapi2.Client, tokenEnsurer state.TokenEnsurer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("remove-label [FLAGS] %s LABEL", strings.ToUpper(lc.ResourceNameSingular)),
		Short: lc.ShortDescriptionRemove,
		Args:  cobra.RangeArgs(1, 2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(lc.NameSuggestions(client)),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) != 1 {
					return nil
				}
				return lc.LabelKeySuggestions(client)(args[0])
			}),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateRemoveLabel, tokenEnsurer.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.RunRemove(ctx, client, cmd, args)
		},
	}
	cmd.Flags().BoolP("all", "a", false, "Remove all labels")
	return cmd
}

// RunRemove executes a remove label command
func (lc *LabelCmds) RunRemove(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	idOrName := args[0]

	labels, id, err := lc.FetchLabels(ctx, client, idOrName)
	if err != nil {
		return err
	}

	if all {
		labels = make(map[string]string)
	} else {
		key := args[1]
		if _, ok := labels[key]; !ok {
			return fmt.Errorf("label %s on %s %d does not exist", key, lc.ResourceNameSingular, id)
		}
		delete(labels, key)
	}

	if err := lc.SetLabels(ctx, client, id, labels); err != nil {
		return err
	}

	if all {
		fmt.Printf("All labels removed from %s %d\n", lc.ResourceNameSingular, id)
	} else {
		fmt.Printf("Label %s removed from %s %d\n", args[1], lc.ResourceNameSingular, id)
	}

	return nil
}

func validateRemoveLabel(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	if all && len(args) == 2 {
		return errors.New("must not specify a label key when using --all/-a")
	}
	if !all && len(args) != 2 {
		return errors.New("must specify a label key when not using --all/-a")
	}

	return nil
}
