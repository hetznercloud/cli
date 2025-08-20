package base

import (
	"errors"
	"fmt"
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
type DeleteCmd struct {
	ResourceNameSingular string // e.g. "Server"
	ResourceNamePlural   string // e.g. "Servers"
	ShortDescription     string
	NameSuggestions      func(client hcapi2.Client) func() []string
	AdditionalFlags      func(*cobra.Command)
	Fetch                func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error)
	Delete               func(s state.State, cmd *cobra.Command, resource interface{}) (*hcloud.Action, error)

	// ExperimentalF is a function that will be used to mark the command as experimental.
	ExperimentalF func(state.State, *cobra.Command) *cobra.Command
}

// CobraCommand creates a command that can be registered with cobra.
func (dc *DeleteCmd) CobraCommand(s state.State) *cobra.Command {
	opts := ""
	if dc.AdditionalFlags != nil {
		opts = "[options] "
	}

	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("delete %s<%s>...", opts, util.ToKebabCase(dc.ResourceNameSingular)),
		Short:                 dc.ShortDescription,
		Args:                  util.Validate,
		ValidArgsFunction:     cmpl.SuggestCandidatesF(dc.NameSuggestions(s.Client())),
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
	if dc.ExperimentalF != nil {
		cmd = dc.ExperimentalF(s, cmd)
	}
	return cmd
}

// deleteBatchSize is the batch size when deleting multiple resources in parallel.
const deleteBatchSize = 10

// Run executes a delete command.
func (dc *DeleteCmd) Run(s state.State, cmd *cobra.Command, args []string) error {
	errs := make([]error, 0, len(args))
	deleted := make([]string, 0, len(args))

	for batch := range slices.Chunk(args, deleteBatchSize) {
		results := make([]util.ResourceState, len(batch))
		actions := make([]*hcloud.Action, 0, len(batch))

		for i, idOrName := range batch {
			results[i] = util.ResourceState{IDOrName: idOrName}

			resource, _, err := dc.Fetch(s, cmd, idOrName)
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
