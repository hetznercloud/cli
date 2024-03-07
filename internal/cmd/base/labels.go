package base

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

// LabelCmds allows defining commands for adding labels to resources.
type LabelCmds struct {
	ResourceNameSingular   string
	ShortDescriptionAdd    string
	ShortDescriptionRemove string
	NameSuggestions        func(client hcapi2.Client) func() []string
	LabelKeySuggestions    func(client hcapi2.Client) func(idOrName string) []string
	FetchLabels            func(s state.State, idOrName string) (map[string]string, int64, error)
	SetLabels              func(s state.State, id int64, labels map[string]string) error
}

// AddCobraCommand creates a command that can be registered with cobra.
func (lc *LabelCmds) AddCobraCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("add-label [--overwrite] <%s> <label>...", strings.ToLower(lc.ResourceNameSingular)),
		Short:                 lc.ShortDescriptionAdd,
		Args:                  cobra.MinimumNArgs(2),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(lc.NameSuggestions(s.Client()))),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateAddLabel, s.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.RunAdd(s, cmd, args)
		},
	}
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

// RunAdd executes an add label command
func (lc *LabelCmds) RunAdd(s state.State, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	idOrName := args[0]

	labels, id, err := lc.FetchLabels(s, idOrName)
	if err != nil {
		return err
	}

	if labels == nil {
		labels = map[string]string{}
	}

	var keys []string
	for _, label := range args[1:] {
		key, val := util.SplitLabelVars(label)
		keys = append(keys, key)

		if _, ok := labels[key]; ok && !overwrite {
			return fmt.Errorf("label %s on %s %d already exists", key, lc.ResourceNameSingular, id)
		}

		labels[key] = val
	}

	if err := lc.SetLabels(s, id, labels); err != nil {
		return err
	}

	cmd.Printf("Label(s) %s added to %s %d\n", strings.Join(keys, ", "), lc.ResourceNameSingular, id)
	return nil
}

func validateAddLabel(_ *cobra.Command, args []string) error {
	for _, label := range args[1:] {
		if len(util.SplitLabel(label)) != 2 {
			return fmt.Errorf("invalid label: %s", label)
		}
	}

	return nil
}

// RemoveCobraCommand creates a command that can be registered with cobra.
func (lc *LabelCmds) RemoveCobraCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("remove-label <%s> (--all | <label>...)", strings.ToLower(lc.ResourceNameSingular)),
		Short: lc.ShortDescriptionRemove,
		Args:  cobra.MinimumNArgs(1),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(lc.NameSuggestions(s.Client())),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) != 1 {
					return nil
				}
				return lc.LabelKeySuggestions(s.Client())(args[0])
			}),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateRemoveLabel, s.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.RunRemove(s, cmd, args)
		},
	}
	cmd.Flags().BoolP("all", "a", false, "Remove all labels")
	return cmd
}

// RunRemove executes a remove label command
func (lc *LabelCmds) RunRemove(s state.State, cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	idOrName := args[0]

	labels, id, err := lc.FetchLabels(s, idOrName)
	if err != nil {
		return err
	}

	if all {
		labels = make(map[string]string)
	} else {
		for _, key := range args[1:] {
			if _, ok := labels[key]; !ok {
				return fmt.Errorf("label %s on %s %d does not exist", key, lc.ResourceNameSingular, id)
			}
			delete(labels, key)
		}
	}

	if err := lc.SetLabels(s, id, labels); err != nil {
		return err
	}

	if all {
		cmd.Printf("All labels removed from %s %d\n", lc.ResourceNameSingular, id)
	} else {
		cmd.Printf("Label(s) %s removed from %s %d\n", strings.Join(args[1:], ", "), lc.ResourceNameSingular, id)
	}

	return nil
}

func validateRemoveLabel(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	if all && len(args) > 1 {
		return errors.New("must not specify a label key when using --all/-a")
	}
	if !all && len(args) <= 1 {
		return errors.New("must specify a label key when not using --all/-a")
	}

	return nil
}
