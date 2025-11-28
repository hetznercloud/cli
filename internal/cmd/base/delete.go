package base

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// DeleteCmd allows defining commands for deleting a resource.
type DeleteCmd[T any] struct {
	ResourceNameSingular string // e.g. "Server"
	ResourceNamePlural   string // e.g. "Servers"
	ShortDescription     string
	NameSuggestions      func(client hcapi2.Client) func() []string
	AdditionalFlags      func(*cobra.Command)
	Fetch                FetchFunc[T]
	Delete               func(s state.State, cmd *cobra.Command, resource T) (*hcloud.Action, error)

	// FetchFunc is a factory function that produces [DeleteCmd.Fetch]. Should be set in case the resource has
	// more than a single identifier that is used in the positional arguments.
	// See [DeleteCmd.PositionalArgumentOverride].
	FetchFunc func(s state.State, cmd *cobra.Command, args []string) (FetchFunc[T], error)

	// In case the resource does not have a single identifier that matches [DeleteCmd.ResourceNameSingular], this field
	// can be set to define the list of positional arguments.
	// For example, passing:
	//     []string{"a", "b", "c"}
	// Would result in the usage string:
	//     <a> <b> <c>...
	// Where c is are resources to be deleted.
	PositionalArgumentOverride []string

	// Can be set if the default [DeleteCmd.NameSuggestions] is not enough. This is usually the case when
	// [DeleteCmd.FetchWithArgs] and [DeleteCmd.PositionalArgumentOverride] is being used.
	ValidArgsFunction func(client hcapi2.Client) []cobra.CompletionFunc

	// Experimental is a function that will be used to mark the command as experimental.
	Experimental func(state.State, *cobra.Command) *cobra.Command
}

type FetchFunc[T any] func(s state.State, cmd *cobra.Command, idOrName string) (T, *hcloud.Response, error)

// CobraCommand creates a command that can be registered with cobra.
func (dc *DeleteCmd[T]) CobraCommand(s state.State) *cobra.Command {
	var suggestArgs []cobra.CompletionFunc
	switch {
	case dc.NameSuggestions != nil:
		suggestArgs = append(suggestArgs,
			cmpl.SuggestCandidatesF(dc.NameSuggestions(s.Client())),
		)
	case dc.ValidArgsFunction != nil:
		suggestArgs = append(suggestArgs, dc.ValidArgsFunction(s.Client())...)
	default:
		log.Fatalf("delete command %s is missing ValidArgsFunction or NameSuggestions", dc.ResourceNameSingular)
	}

	opts := ""
	if dc.AdditionalFlags != nil {
		opts = "[options] "
	}

	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("delete %s%s...", opts, positionalArguments(dc.ResourceNameSingular, dc.PositionalArgumentOverride)),
		Short:                 dc.ShortDescription,
		Args:                  util.Validate,
		ValidArgsFunction:     cmpl.SuggestArgs(suggestArgs...),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(s.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dc.Run(s, cmd, args)
		},
	}
	if dc.AdditionalFlags != nil {
		dc.AdditionalFlags(cmd)
	}
	if dc.Experimental != nil {
		cmd = dc.Experimental(s, cmd)
	}
	return cmd
}

// deleteBatchSize is the batch size when deleting multiple resources in parallel.
const deleteBatchSize = 10

// Run executes a delete command.
func (dc *DeleteCmd[T]) Run(s state.State, cmd *cobra.Command, args []string) error {
	toDelete := args[max(0, len(dc.PositionalArgumentOverride)-1):]

	errs := make([]error, 0, len(toDelete))
	deleted := make([]string, 0, len(toDelete))

	var (
		fetch FetchFunc[T]
		err   error
	)
	if dc.FetchFunc != nil {
		fetch, err = dc.FetchFunc(s, cmd, args)
		if err != nil {
			return err
		}
	} else {
		fetch = dc.Fetch
	}

	for batch := range slices.Chunk(toDelete, deleteBatchSize) {
		results := make([]util.ResourceState, len(batch))
		actions := make([]*hcloud.Action, 0, len(batch))

		for i, idOrName := range batch {
			results[i] = util.ResourceState{IDOrName: idOrName}

			resource, _, err := fetch(s, cmd, idOrName)
			if err != nil {
				results[i].Error = err
				continue
			}
			if util.IsNil(resource) {
				results[i].Error = fmt.Errorf("%s not found: %s", dc.ResourceNameSingular, idOrName)
				continue
			}

			action, err := dc.Delete(s, cmd, resource)
			if err != nil {
				results[i].Error = err
				continue
			}
			if action != nil {
				actions = append(actions, action)
			}
		}

		for _, result := range results {
			if result.Error != nil {
				errs = append(errs, result.Error)
			} else {
				deleted = append(deleted, result.IDOrName)
			}
		}

		if len(actions) > 0 {
			// TODO: We do not check if an action fails for a specific resource
			if err := s.WaitForActions(s, cmd, actions...); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(deleted) == 1 {
		cmd.Printf("%s %s deleted\n", dc.ResourceNameSingular, deleted[0])
	} else if len(deleted) > 1 {
		cmd.Printf("%s %s deleted\n", dc.ResourceNamePlural, strings.Join(deleted, ", "))
	}

	return errors.Join(errs...)
}
