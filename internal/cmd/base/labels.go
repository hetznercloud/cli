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
	Fetch                  func(s state.State, idOrName string) (any, error)
	SetLabels              func(s state.State, resource any, labels map[string]string) error
	GetLabels              func(resource any) map[string]string
	GetIDOrName            func(resource any) string
}

// AddCobraCommand creates a command that can be registered with cobra.
func (lc *LabelCmds) AddCobraCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("add-label [--overwrite] <%s> <label>...", util.ToKebabCase(lc.ResourceNameSingular)),
		Short:                 lc.ShortDescriptionAdd,
		Args:                  util.Validate,
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

	resource, err := lc.Fetch(s, idOrName)
	if err != nil {
		return err
	}

	labels, idOrName := lc.GetLabels(resource), lc.GetIDOrName(resource)

	if labels == nil {
		labels = map[string]string{}
	}

	var keys []string
	for _, label := range args[1:] {
		key, val := util.SplitLabelVars(label)
		keys = append(keys, key)

		if _, ok := labels[key]; ok && !overwrite {
			return fmt.Errorf("label %s on %s %s already exists", key, lc.ResourceNameSingular, idOrName)
		}

		labels[key] = val
	}

	if err := lc.SetLabels(s, resource, labels); err != nil {
		return err
	}

	cmd.Printf("Label(s) %s added to %s %s\n", strings.Join(keys, ", "), lc.ResourceNameSingular, idOrName)
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
		Use:   fmt.Sprintf("remove-label <%s> (--all | <label>...)", util.ToKebabCase(lc.ResourceNameSingular)),
		Short: lc.ShortDescriptionRemove,
		Args:  util.ValidateLenient,
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

	resource, err := lc.Fetch(s, idOrName)
	if err != nil {
		return err
	}

	labels, idOrName := lc.GetLabels(resource), lc.GetIDOrName(resource)

	if all {
		labels = make(map[string]string)
	} else {
		for _, key := range args[1:] {
			if _, ok := labels[key]; !ok {
				return fmt.Errorf("label %s on %s %s does not exist", key, lc.ResourceNameSingular, idOrName)
			}
			delete(labels, key)
		}
	}

	if err := lc.SetLabels(s, resource, labels); err != nil {
		return err
	}

	if all {
		cmd.Printf("All labels removed from %s %s\n", lc.ResourceNameSingular, idOrName)
	} else {
		cmd.Printf("Label(s) %s removed from %s %s\n", strings.Join(args[1:], ", "), lc.ResourceNameSingular, idOrName)
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
