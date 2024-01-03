package server

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func getChangeProtectionOpts(enable bool, flags []string) (hcloud.ServerChangeProtectionOpts, error) {

	opts := hcloud.ServerChangeProtectionOpts{}

	var unknown []string
	for _, arg := range flags {
		switch strings.ToLower(arg) {
		case "delete":
			opts.Delete = hcloud.Ptr(enable)
		case "rebuild":
			opts.Rebuild = hcloud.Ptr(enable)
		default:
			unknown = append(unknown, arg)
		}
	}
	if len(unknown) > 0 {
		return opts, fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
	}

	return opts, nil
}

func changeProtection(s state.State, cmd *cobra.Command,
	server *hcloud.Server, enable bool, opts hcloud.ServerChangeProtectionOpts) error {

	if opts.Delete == nil && opts.Rebuild == nil {
		return nil
	}

	action, _, err := s.Server().ChangeProtection(s, server, opts)
	if err != nil {
		return err
	}

	if err := s.ActionProgress(cmd, s, action); err != nil {
		return err
	}

	if enable {
		cmd.Printf("Resource protection enabled for server %d\n", server.ID)
	} else {
		cmd.Printf("Resource protection disabled for server %d\n", server.ID)
	}
	return nil
}

var EnableProtectionCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "enable-protection [FLAGS] SERVER PROTECTIONLEVEL [PROTECTIONLEVEL...]",
			Short: "Enable resource protection for a server",
			Args:  cobra.MinimumNArgs(2),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.Server().Names),
				cmpl.SuggestCandidates("delete", "rebuild"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := s.Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		opts, err := getChangeProtectionOpts(true, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(s, cmd, server, true, opts)
	},
}
