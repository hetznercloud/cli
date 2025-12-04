package base

import (
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// ChangeProtectionCmds allows defining commands for changing a resource's protection.
type ChangeProtectionCmds[T, Opts any] struct {
	ResourceNameSingular    string // e.g. "Server"
	ShortEnableDescription  string
	ShortDisableDescription string
	NameSuggestions         func(client hcapi2.Client) func() []string
	AdditionalFlags         func(*cobra.Command)
	// Fetch is called to fetch the resource to describe.
	// The first returned interface is the resource itself as a hcloud struct, the second is the schema for the resource.
	Fetch func(s state.State, cmd *cobra.Command, idOrName string) (T, *hcloud.Response, error)
	// Can be set in case the resource has more than a single identifier that is used in the positional arguments.
	// See [DescribeCmd.PositionalArgumentOverride].
	FetchWithArgs func(s state.State, cmd *cobra.Command, args []string) (T, *hcloud.Response, error)

	// In case the resource does not have a single identifier that matches [DescribeCmd.ResourceNameSingular], this field
	// can be set to define the list of positional arguments.
	// For example, passing:
	//     []string{"a", "b", "c"}
	// Would result in the usage string:
	//     <a> <b> <c>
	PositionalArgumentOverride []string

	// Can be set if the default [DescribeCmd.NameSuggestions] is not enough. This is usually the case when
	// [DescribeCmd.FetchWithArgs] and [DescribeCmd.PositionalArgumentOverride] is being used.
	ValidArgsFunction func(client hcapi2.Client) []cobra.CompletionFunc

	// Levels maps all available protection levels to a function that sets the corresponding value in the Opts struct
	ProtectionLevels map[string]func(opts *Opts, value bool)

	// If ProtectionLevelOptional is set, all protection levels will always be applied
	ProtectionLevelsOptional bool

	// ChangeProtectionFunction is used to change the protection on a specific resource given the Opts
	ChangeProtectionFunction func(s state.State, resource T, opts Opts) (*hcloud.Action, *hcloud.Response, error)

	// GetIDFunction is used to retrieve the ID of a resource
	IDOrName func(resource T) string

	// Experimental is a function that will be used to mark the command as experimental.
	Experimental func(state.State, *cobra.Command) *cobra.Command
}

func (cpc *ChangeProtectionCmds[T, Opts]) newChangeProtectionCmd(s state.State, enable bool) *cobra.Command {
	if len(cpc.ProtectionLevels) < 1 {
		log.Fatalf("change protection command %s is missing ProtectionLevels", cpc.ResourceNameSingular)
	}

	levels := maps.Keys(cpc.ProtectionLevels)
	slices.Sort(levels)

	var suggestArgs []cobra.CompletionFunc
	switch {
	case cpc.NameSuggestions != nil:
		suggestArgs = append(suggestArgs,
			cmpl.SuggestCandidatesF(cpc.NameSuggestions(s.Client())),
		)
	case cpc.ValidArgsFunction != nil:
		suggestArgs = append(suggestArgs, cpc.ValidArgsFunction(s.Client())...)
	default:
		log.Fatalf("change protection command %s is missing ValidArgsFunction or NameSuggestions", cpc.ResourceNameSingular)
	}

	suggestArgs = append(suggestArgs, cmpl.SuggestCandidates(levels...))

	var cmdName string
	if enable {
		cmdName = "enable-protection"
	} else {
		cmdName = "disable-protection"
	}

	var levelUsage string
	if len(levels) == 1 {
		if cpc.ProtectionLevelsOptional {
			levelUsage = fmt.Sprintf("[%s]", levels[0])
		} else {
			levelUsage = levels[0]
		}
	} else {
		if cpc.ProtectionLevelsOptional {
			levelUsage = fmt.Sprintf("[%s]...", strings.Join(levels, "|"))
		} else {
			levelUsage = fmt.Sprintf("(%s)...", strings.Join(levels, "|"))
		}
	}

	var shortDescription string
	if enable {
		if cpc.ShortEnableDescription != "" {
			shortDescription = cpc.ShortEnableDescription
		} else {
			shortDescription = fmt.Sprintf("Enable resource protection for a %s", cpc.ResourceNameSingular)
		}
	} else {
		if cpc.ShortDisableDescription != "" {
			shortDescription = cpc.ShortDisableDescription
		} else {
			shortDescription = fmt.Sprintf("Disable resource protection for a %s", cpc.ResourceNameSingular)
		}
	}

	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("%s %s %s", cmdName, positionalArguments(cpc.ResourceNameSingular, cpc.PositionalArgumentOverride), levelUsage),
		Short:                 shortDescription,
		Args:                  util.ValidateLenient,
		ValidArgsFunction:     cmpl.SuggestArgs(suggestArgs...),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(s.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cpc.Run(s, cmd, args, enable)
		},
	}

	if cpc.AdditionalFlags != nil {
		cpc.AdditionalFlags(cmd)
	}

	if cpc.Experimental != nil {
		cmd = cpc.Experimental(s, cmd)
	}

	return cmd
}

// EnableCobraCommand creates an enable-protection command that can be registered with cobra.
func (cpc *ChangeProtectionCmds[T, Opts]) EnableCobraCommand(s state.State) *cobra.Command {
	return cpc.newChangeProtectionCmd(s, true)
}

// DisableCobraCommand creates a disable-protection command that can be registered with cobra.
func (cpc *ChangeProtectionCmds[T, Opts]) DisableCobraCommand(s state.State) *cobra.Command {
	return cpc.newChangeProtectionCmd(s, false)
}

// Run executes a describe command.
func (cpc *ChangeProtectionCmds[T, Opts]) Run(s state.State, cmd *cobra.Command, args []string, enable bool) error {

	var (
		resource T
		err      error
	)
	if cpc.FetchWithArgs != nil {
		resource, _, err = cpc.FetchWithArgs(s, cmd, args)
	} else {
		resource, _, err = cpc.Fetch(s, cmd, args[0])
	}
	if err != nil {
		return err
	}

	// resource is an interface that always has a type, so the interface is never nil
	// (i.e. == nil) is always false.
	if reflect.ValueOf(resource).IsNil() {
		args := args[:max(1, len(cpc.PositionalArgumentOverride))]
		return fmt.Errorf("%s not found: %s", cpc.ResourceNameSingular, strings.Join(args, " "))
	}

	levels := args[max(1, len(cpc.PositionalArgumentOverride)):]
	if cpc.ProtectionLevelsOptional {
		levels = maps.Keys(cpc.ProtectionLevels)
	}

	opts, err := cpc.GetChangeProtectionOpts(enable, levels)
	if err != nil {
		return err
	}

	return cpc.ChangeProtection(s, cmd, resource, enable, opts)
}

func (cpc *ChangeProtectionCmds[T, Opts]) GetChangeProtectionOpts(enable bool, levels []string) (Opts, error) {
	var (
		opts    Opts
		unknown []string
	)
	for _, level := range levels {
		if f, ok := cpc.ProtectionLevels[strings.ToLower(level)]; ok {
			f(&opts, enable)
		} else {
			unknown = append(unknown, level)
		}
	}
	if len(unknown) > 0 {
		return opts, fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
	}
	return opts, nil
}

func (cpc *ChangeProtectionCmds[T, Opts]) ChangeProtection(s state.State, cmd *cobra.Command,
	resource T, enable bool, opts Opts) error {

	action, _, err := cpc.ChangeProtectionFunction(s, resource, opts)
	if err != nil {
		return err
	}

	if err := s.WaitForActions(s, cmd, action); err != nil {
		return err
	}

	idOrName := cpc.IDOrName(resource)
	if enable {
		cmd.Printf("Resource protection enabled for %s %s\n", cpc.ResourceNameSingular, idOrName)
	} else {
		cmd.Printf("Resource protection disabled for %s %s\n", cpc.ResourceNameSingular, idOrName)
	}
	return nil
}
