package base

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// SetRdnsCmd allows defining commands for setting the RDNS of a resource.
type SetRdnsCmd struct {
	ResourceNameSingular string // e.g. "server"
	ShortDescription     string
	NameSuggestions      func(client hcapi2.Client) func() []string
	AdditionalFlags      func(*cobra.Command)
	Fetch                func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error)
	GetDefaultIP         func(resource interface{}) net.IP
}

// CobraCommand creates a command that can be registered with cobra.
func (rc *SetRdnsCmd) CobraCommand(
	ctx context.Context, client hcapi2.Client, tokenEnsurer state.TokenEnsurer, actionWaiter state.ActionWaiter,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("set-rdns [FLAGS] %s", strings.ToUpper(rc.ResourceNameSingular)),
		Short:                 rc.ShortDescription,
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(rc.NameSuggestions(client))),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(tokenEnsurer.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return rc.Run(ctx, client, cmd, actionWaiter, args)
		},
	}
	cmd.Flags().StringP("hostname", "r", "", "Hostname to set as a reverse DNS PTR entry (required)")
	cmd.MarkFlagRequired("hostname")

	cmd.Flags().IPP("ip", "i", net.IP{}, "IP address for which the reverse DNS entry should be set")
	if rc.AdditionalFlags != nil {
		rc.AdditionalFlags(cmd)
	}
	return cmd
}

// Run executes a setRDNS command.
func (rc *SetRdnsCmd) Run(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, actionWaiter state.ActionWaiter, args []string) error {

	idOrName := args[0]
	resource, _, err := rc.Fetch(ctx, client, cmd, idOrName)
	if err != nil {
		return err
	}

	// resource is an interface that always has a type, so the interface is never nil
	// (i.e. == nil) is always false.
	if reflect.ValueOf(resource).IsNil() {
		return fmt.Errorf("%s not found: %s", rc.ResourceNameSingular, idOrName)
	}

	ip, _ := cmd.Flags().GetIP("ip")
	if ip.IsUnspecified() || ip == nil {
		ip = rc.GetDefaultIP(resource)
	}
	hostname, _ := cmd.Flags().GetString("hostname")
	action, _, err := client.RDNS().ChangeDNSPtr(ctx, resource.(hcloud.RDNSSupporter), ip, hcloud.String(hostname))
	if err != nil {
		return err
	}

	if err := actionWaiter.ActionProgress(ctx, action); err != nil {
		return err
	}

	fmt.Printf("Reverse DNS of %s %s changed\n", rc.ResourceNameSingular, idOrName)

	return nil
}
