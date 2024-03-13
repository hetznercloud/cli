package base

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// UpdateCmd allows defining commands for updating a resource.
type UpdateCmd struct {
	ResourceNameSingular string // e.g. "server"
	ShortDescription     string
	NameSuggestions      func(client hcapi2.Client) func() []string
	DefineFlags          func(*cobra.Command)
	Fetch                func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error)
	Update               func(s state.State, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error
}

// CobraCommand creates a command that can be registered with cobra.
func (uc *UpdateCmd) CobraCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("update [options] <%s>", strings.ToLower(uc.ResourceNameSingular)),
		Short:                 uc.ShortDescription,
		Args:                  util.Validate,
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(uc.NameSuggestions(s.Client()))),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(s.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return uc.Run(s, cmd, args)
		},
	}
	uc.DefineFlags(cmd)
	return cmd
}

// Run executes a update command.
func (uc *UpdateCmd) Run(s state.State, cmd *cobra.Command, args []string) error {

	idOrName := args[0]
	resource, _, err := uc.Fetch(s, cmd, idOrName)
	if err != nil {
		return err
	}

	// resource is an interface that always has a type, so the interface is never nil
	// (i.e. == nil) is always false.
	if reflect.ValueOf(resource).IsNil() {
		return fmt.Errorf("%s not found: %s", uc.ResourceNameSingular, idOrName)
	}

	// The inherited commands should not need to parse the flags themselves
	// or use the cobra command, therefore we fill them in a map here and
	// pass the map then to the update method. A caller can/should rely on
	// the map to contain all the flag keys that were specified.
	flags := make(map[string]pflag.Value, cmd.Flags().NFlag())
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flags[flag.Name] = flag.Value
	})

	if err := uc.Update(s, cmd, resource, flags); err != nil {
		return fmt.Errorf("updating %s %s failed: %s", uc.ResourceNameSingular, idOrName, err)
	}

	cmd.Printf("%s %v updated\n", uc.ResourceNameSingular, idOrName)
	return nil
}
