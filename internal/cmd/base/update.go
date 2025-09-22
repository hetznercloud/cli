package base

import (
	"fmt"
	"log"
	"reflect"

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
	ResourceNameSingular string // e.g. "Server"
	ShortDescription     string
	NameSuggestions      func(client hcapi2.Client) func() []string
	DefineFlags          func(*cobra.Command)

	Fetch func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error)
	// Can be set in case the resource has more than a single identifier that is used in the positional arguments.
	// See [UpdateCmd.PositionalArgumentOverride].
	FetchWithArgs func(s state.State, cmd *cobra.Command, args []string) (any, *hcloud.Response, error)

	Update func(s state.State, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error

	// In case the resource does not have a single identifier that matches [UpdateCmd.ResourceNameSingular], this field
	// can be set to define the list of positional arguments.
	// For example, passing:
	//     []string{"a", "b", "c"}
	// Would result in the usage string:
	//     <a> <b> <c>
	PositionalArgumentOverride []string

	// Can be set if the default [UpdateCmd.NameSuggestions] is not enough. This is usually the case when
	// [UpdateCmd.FetchWithArgs] and [UpdateCmd.PositionalArgumentOverride] is being used.
	ValidArgsFunction func(client hcapi2.Client) []cobra.CompletionFunc

	// Experimental is a function that will be used to mark the command as experimental.
	Experimental func(state.State, *cobra.Command) *cobra.Command
}

// CobraCommand creates a command that can be registered with cobra.
func (uc *UpdateCmd) CobraCommand(s state.State) *cobra.Command {
	var suggestArgs []cobra.CompletionFunc
	switch {
	case uc.NameSuggestions != nil:
		suggestArgs = append(suggestArgs,
			cmpl.SuggestCandidatesF(uc.NameSuggestions(s.Client())),
		)
	case uc.ValidArgsFunction != nil:
		suggestArgs = append(suggestArgs, uc.ValidArgsFunction(s.Client())...)
	default:
		log.Fatalf("update command %s is missing ValidArgsFunction or NameSuggestions", uc.ResourceNameSingular)
	}

	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("update [options] %s", positionalArguments(uc.ResourceNameSingular, uc.PositionalArgumentOverride)),
		Short:                 uc.ShortDescription,
		Args:                  util.Validate,
		ValidArgsFunction:     cmpl.SuggestArgs(suggestArgs...),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(s.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return uc.Run(s, cmd, args)
		},
	}
	uc.DefineFlags(cmd)
	if uc.Experimental != nil {
		cmd = uc.Experimental(s, cmd)
	}
	return cmd
}

// Run executes a update command.
func (uc *UpdateCmd) Run(s state.State, cmd *cobra.Command, args []string) error {
	var (
		resource any
		err      error
	)
	if uc.FetchWithArgs != nil {
		resource, _, err = uc.FetchWithArgs(s, cmd, args)
	} else {
		resource, _, err = uc.Fetch(s, cmd, args[0])
	}
	if err != nil {
		return err
	}

	idOrName := args[len(args)-1]

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
		return fmt.Errorf("updating %s %s failed: %w", uc.ResourceNameSingular, idOrName, err)
	}

	cmd.Printf("%s %v updated\n", uc.ResourceNameSingular, idOrName)
	return nil
}
