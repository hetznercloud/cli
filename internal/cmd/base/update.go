package base

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

// UpdateCmd allows defining commands for updating a resource.
type UpdateCmd struct {
	ResourceNameSingular string // e.g. "server"
	ShortDescription     string
	NameSuggestions      func(client hcapi2.Client) func() []string
	DefineFlags          func(*cobra.Command)
	Fetch                func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error)
	Update               func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error
}

// CobraCommand creates a command that can be registered with cobra.
func (uc *UpdateCmd) CobraCommand(
	ctx context.Context, client hcapi2.Client, tokenEnsurer state.TokenEnsurer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("update [FLAGS] %s", strings.ToUpper(uc.ResourceNameSingular)),
		Short:                 uc.ShortDescription,
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(uc.NameSuggestions(client))),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(tokenEnsurer.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return uc.Run(ctx, client, cmd, args)
		},
	}
	uc.DefineFlags(cmd)
	return cmd
}

// Run executes a update command.
func (uc *UpdateCmd) Run(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, args []string) error {

	idOrName := args[0]
	resource, _, err := uc.Fetch(ctx, client, cmd, idOrName)
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

	if err := uc.Update(ctx, client, cmd, resource, flags); err != nil {
		return fmt.Errorf("updating %s %s failed: %s", uc.ResourceNameSingular, idOrName, err)
	}

	fmt.Printf("%s %v updated\n", uc.ResourceNameSingular, idOrName)
	return nil
}
