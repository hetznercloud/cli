package base

import (
	"errors"
	"fmt"
	"reflect"

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
	ShortDescription     string
	NameSuggestions      func(client hcapi2.Client) func() []string
	AdditionalFlags      func(*cobra.Command)
	Fetch                func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error)
	Delete               func(s state.State, cmd *cobra.Command, resource interface{}) error
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
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(dc.NameSuggestions(s.Client()))),
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

// Run executes a describe command.
func (dc *DeleteCmd) Run(s state.State, cmd *cobra.Command, args []string) error {
	var cmdErr error

	for _, idOrName := range args {
		resource, _, err := dc.Fetch(s, cmd, idOrName)
		if err != nil {
			cmdErr = errors.Join(cmdErr, err)
			continue
		}

		// resource is an interface that always has a type, so the interface is never nil
		// (i.e. == nil) is always false.
		if reflect.ValueOf(resource).IsNil() {
			cmdErr = errors.Join(cmdErr, fmt.Errorf("%s not found: %s", dc.ResourceNameSingular, idOrName))
			continue
		}

		if err = dc.Delete(s, cmd, resource); err != nil {
			cmdErr = errors.Join(cmdErr, fmt.Errorf("deleting %s %s failed: %s", dc.ResourceNameSingular, idOrName, err))
		}
		cmd.Printf("%s %v deleted\n", dc.ResourceNameSingular, idOrName)
	}

	return cmdErr
}
