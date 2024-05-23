package base

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// DeleteCmd allows defining commands for deleting a resource.
type DeleteCmd struct {
	ResourceNameSingular string // e.g. "server"
	ResourceNamePlural   string // e.g. "servers"
	ShortDescription     string
	NameSuggestions      func(client hcapi2.Client) func() []string
	AdditionalFlags      func(*cobra.Command)
	Fetch                func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error)
	Delete               func(s state.State, cmd *cobra.Command, resource interface{}) (*hcloud.Action, error)
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
	return cmd
}

// Run executes a delete command.
func (dc *DeleteCmd) Run(s state.State, cmd *cobra.Command, args []string) error {

	wg := sync.WaitGroup{}
	wg.Add(len(args))
	actions, errs :=
		make([]*hcloud.Action, len(args)),
		make([]error, len(args))

	for i, idOrName := range args {
		i, idOrName := i, idOrName
		go func() {
			defer wg.Done()
			resource, _, err := dc.Fetch(s, cmd, idOrName)
			if err != nil {
				errs[i] = err
				return
			}
			if util.IsNil(resource) {
				errs[i] = fmt.Errorf("%s not found: %s", dc.ResourceNameSingular, idOrName)
				return
			}
			actions[i], errs[i] = dc.Delete(s, cmd, resource)
		}()
	}

	wg.Wait()
	filtered := util.FilterNil(actions)
	var err error
	if len(filtered) > 0 {
		err = s.WaitForActions(cmd, s, util.FilterNil(actions)...)
	}

	var actuallyDeleted []string
	for i, idOrName := range args {
		if errs[i] == nil {
			actuallyDeleted = append(actuallyDeleted, idOrName)
		}
	}

	if len(actuallyDeleted) == 1 {
		cmd.Printf("%s %s deleted\n", dc.ResourceNameSingular, actuallyDeleted[0])
	} else if len(actuallyDeleted) > 1 {
		cmd.Printf("%s %s deleted\n", dc.ResourceNamePlural, strings.Join(actuallyDeleted, ", "))
	}
	return errors.Join(append(errs, err)...)
}
