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

// RemoveLabelCmd allows defining commands for adding labels to resources.
type RemoveLabelCmd struct {
	ResourceNameSingular string
	ShortDescription     string
	NameSuggestions      func(client hcapi2.Client) func() []string
	LabelKeySuggestions  func(client hcapi2.Client) func(idOrName string) []string
	FetchLabels          func(ctx context.Context, client hcapi2.Client, idOrName string) (map[string]string, int, error)
	SetLabels            func(ctx context.Context, client hcapi2.Client, id int, labels map[string]string) error
}

// CobraCommand creates a command that can be registered with cobra.
func (lc *RemoveLabelCmd) CobraCommand(
	ctx context.Context, client hcapi2.Client, tokenEnsurer state.TokenEnsurer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("remove-label [FLAGS] %s LABEL", strings.ToUpper(lc.ResourceNameSingular)),
		Short: lc.ShortDescription,
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
			return lc.Run(ctx, client, cmd, args)
		},
	}
	cmd.Flags().BoolP("all", "a", false, "Remove all labels")
	return cmd
}

// Run executes a remove label command
func (lc *RemoveLabelCmd) Run(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, args []string) error {
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
