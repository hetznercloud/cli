package context

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func newSSHKeyCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ssh-key",
		Short:                 "Manage a context's default SSH key",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	cmd.PersistentFlags().String("context", "", "Name of the context to manage the default SSH keys for (Default: Active Context)")
	_ = cmd.RegisterFlagCompletionFunc("context", cmpl.SuggestCandidates(config.ContextNames(s.Config())...))

	cmd.AddCommand(newSSHKeyAddCommand(s))
	cmd.AddCommand(newSSHKeyRemoveCommand(s))
	cmd.AddCommand(newSSHKeyListCommand(s))
	return cmd
}

func newSSHKeyAddCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add SSH-KEY...",
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(s.Client().SSHKey().Names)),
		Short:                 "Add a default SSH key to the context",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(s, runSSHKeyAdd),
	}

	cmd.Flags().Bool("all", false, "Add all available SSH keys to the context")
	_ = cmd.RegisterFlagCompletionFunc("all", cmpl.SuggestCandidates("true", "false"))
	return cmd
}

func runSSHKeyAdd(s state.State, cmd *cobra.Command, args []string) error {

	ctx, err := getContext(s, cmd)
	if err != nil {
		return err
	}

	s.Client().WithOpts(hcloud.WithToken(ctx.Token))
	keys := args

	all, _ := cmd.Flags().GetBool("all")
	if all {
		allKeys, err := s.Client().SSHKey().All(s)
		if err != nil {
			return err
		}
		if len(allKeys) == 0 {
			return fmt.Errorf("no SSH keys available")
		}
		keys = make([]string, len(allKeys))
		for i, key := range allKeys {
			keys[i] = fmt.Sprintf("%d", key.ID)
		}
	} else {
		if len(keys) == 0 {
			return fmt.Errorf("no SSH keys specified")
		}
		var (
			notExist []string
			wg       = &sync.WaitGroup{}
			mu       = &sync.Mutex{}
		)
		wg.Add(len(keys))
		for _, key := range keys {
			key := key
			go func() {
				k, _, _ := s.Client().SSHKey().Get(s, key)
				if k == nil {
					mu.Lock()
					notExist = append(notExist, key)
					mu.Unlock()
				}
				wg.Done()
			}()
		}
		wg.Wait()
		if len(notExist) > 0 {
			_, _ = fmt.Fprintf(os.Stderr, "Warning: The given SSH key(s) %s do not exist in context \"%s\"\n", strings.Join(notExist, ", "), ctx.Name)
		}
	}

	ctx.SSHKeys = append(ctx.SSHKeys, keys...)
	// remove duplicates
	slices.Sort(ctx.SSHKeys)
	ctx.SSHKeys = slices.Compact(ctx.SSHKeys)

	if err := s.Config().Write(); err != nil {
		return err
	}

	cmd.Printf("Added SSH key(s) %s to context \"%s\"\n", strings.Join(keys, ", "), ctx.Name)
	return nil
}

func newSSHKeyRemoveCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "remove SSH-KEY...",
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidates(config.ContextNames(s.Config())...)),
		Short:                 "Remove a default SSH key from the context",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(s, runSSHKeyRemove),
	}

	cmd.Flags().Bool("all", false, "Remove all SSH keys from the context")
	_ = cmd.RegisterFlagCompletionFunc("all", cmpl.SuggestCandidates("true", "false"))
	return cmd
}

func runSSHKeyRemove(s state.State, cmd *cobra.Command, args []string) error {

	ctx, err := getContext(s, cmd)
	if err != nil {
		return err
	}

	keys := args
	origLen := len(ctx.SSHKeys)

	all, _ := cmd.Flags().GetBool("all")
	if all {
		ctx.SSHKeys = nil
	} else {
		var newKeys []string
		for _, key := range ctx.SSHKeys {
			if slices.Contains(keys, key) {
				continue
			}
			newKeys = append(newKeys, key)
		}
		ctx.SSHKeys = newKeys
	}

	if err := s.Config().Write(); err != nil {
		return err
	}

	removed := origLen - len(ctx.SSHKeys)
	cmd.Printf("Removed %d SSH key(s) from context \"%s\"\n", removed, ctx.Name)
	return nil
}

func newSSHKeyListCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "List all default SSH keys of the context",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(s, runSSHKeyList),
	}
	return cmd
}

func runSSHKeyList(s state.State, cmd *cobra.Command, _ []string) error {

	ctx, err := getContext(s, cmd)
	if err != nil {
		return err
	}

	if len(ctx.SSHKeys) == 0 {
		cmd.Printf("No SSH keys in context \"%s\"\n", ctx.Name)
		return nil
	}

	cmd.Printf("SSH keys in context \"%s\":\n", ctx.Name)
	for _, key := range ctx.SSHKeys {
		cmd.Printf(" - %s\n", key)
	}
	return nil
}

func getContext(s state.State, cmd *cobra.Command) (*config.Context, error) {

	var ctx *config.Context

	ctxName, _ := cmd.Flags().GetString("context")
	if ctxName != "" {
		ctx = config.ContextByName(s.Config(), ctxName)
		if ctx == nil {
			return nil, fmt.Errorf("context not found: %v", ctxName)
		}
	} else {
		ctx = s.Config().ActiveContext()
		if ctx == nil {
			return nil, fmt.Errorf("no active context")
		}
	}

	return ctx, nil
}
