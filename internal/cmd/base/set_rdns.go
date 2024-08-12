package base

import (
	"fmt"
	"net"
	"reflect"

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
	Fetch                func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error)
	GetDefaultIP         func(resource interface{}) net.IP
}

// CobraCommand creates a command that can be registered with cobra.
func (rc *SetRdnsCmd) CobraCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("set-rdns [--ip <ip>] (--hostname <hostname> | --reset) <%s>", util.ToKebabCase(rc.ResourceNameSingular)),
		Short:                 rc.ShortDescription,
		Args:                  util.Validate,
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(rc.NameSuggestions(s.Client()))),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(s.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return rc.Run(s, cmd, args)
		},
	}
	cmd.Flags().StringP("hostname", "r", "", "Hostname to set as a reverse DNS PTR entry")
	cmd.Flags().Bool("reset", false, "Reset the reverse DNS entry to the default value")

	cmd.Flags().IPP("ip", "i", net.IP{}, "IP address for which the reverse DNS entry should be set")
	return cmd
}

// Run executes a setRDNS command.
func (rc *SetRdnsCmd) Run(s state.State, cmd *cobra.Command, args []string) error {

	var hostnamePtr *string
	if reset, _ := cmd.Flags().GetBool("reset"); reset {
		hostnamePtr = nil
	} else {
		hostname, _ := cmd.Flags().GetString("hostname")
		if hostname == "" {
			return fmt.Errorf("either --hostname or --reset must be specified")
		}
		hostnamePtr = &hostname
	}

	idOrName := args[0]
	resource, _, err := rc.Fetch(s, cmd, idOrName)
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

	action, _, err := s.Client().RDNS().ChangeDNSPtr(s, resource.(hcloud.RDNSSupporter), ip, hostnamePtr)
	if err != nil {
		return err
	}

	if err := s.WaitForActions(cmd, s, action); err != nil {
		return err
	}

	cmd.Printf("Reverse DNS of %s %s changed\n", rc.ResourceNameSingular, idOrName)

	return nil
}
