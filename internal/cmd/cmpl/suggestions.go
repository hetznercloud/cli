package cmpl

import (
	"strings"

	"github.com/spf13/cobra"
)

// SuggestCandidates returns a function that selects all items from the list of
// candidates cs which have the prefix toComplete. If toComplete is empty cs is
// returned.
//
// The returned function is mainly intended to be passed to
// cobra/Command.RegisterFlagCompletionFunc or assigned  to
// cobra/Command.ValidArgsFunction.
func SuggestCandidates(cs ...string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return SuggestCandidatesF(func() []string {
		return cs
	})
}

// SuggestCandidatesF returns a function that calls the candidate function cf
// to obtain a list of completion candidates. Once the list of candidates is
// obtained the function returned by SuggestCandidatesF behaves like the
// function returned by SuggestCandidates.
func SuggestCandidatesF(
	cf func() []string,
) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return SuggestCandidatesCtx(func(*cobra.Command, []string) []string {
		return cf()
	})
}

// SuggestCandidatesCtx returns a function that uses the candidate function cf
// to obtain a list of completion candidates in the context of previously
// selected arguments and flags. This is mainly useful for completion candidates that
// depend on a previously selected item like a server.
//
// Once the list of candidates is obtained the function returned by
// SuggestCandidatesF behaves like the function returned by SuggestCandidates.
func SuggestCandidatesCtx(
	cf func(*cobra.Command, []string) []string,
) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		cs := cf(cmd, args)
		if toComplete == "" {
			return cs, cobra.ShellCompDirectiveDefault
		}

		var sel []string
		for _, c := range cs {
			if !strings.HasPrefix(c, toComplete) {
				continue
			}
			sel = append(sel, c)
		}

		return sel, cobra.ShellCompDirectiveDefault
	}

}

// SuggestNothing returns a function that provides no suggestions.
func SuggestNothing() func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveDefault
	}
}

// SuggestArgs returns a function that uses the slice of valid argument
// functions vfs to provide completion suggestions for the passed command line
// arguments.
//
// The selection of the respective entry in vfs is positional, i.e. to
// determine the suggestions for the fourth command line argument SuggestArgs
// calls the function at vfs[4] if it exists. To skip suggestions for an
// argument in the middle of a list of arguments pass either nil or
// SuggestNothing. Using SuggestNothing is preferred.
func SuggestArgs(
	vfs ...func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective),
) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Number of argument to suggest. args contains the already present
		// command line arguments.
		argNo := len(args)

		// Skip completion if not enough vfs have been passed, or if vfs at
		// argNo is nil.
		if len(vfs) <= argNo || vfs[argNo] == nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		f := vfs[argNo]
		return f(cmd, args, toComplete)
	}
}
