package floatingip

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd[*hcloud.FloatingIP]{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "create [options] --type <ipv4|ipv6> (--home-location <location> | --server <server>)",
			Short:                 "Create a Floating IP",
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("type", "", "Type (ipv4 or ipv6) (required)")
		_ = cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("ipv4", "ipv6"))
		_ = cmd.MarkFlagRequired("type")

		cmd.Flags().String("description", "", "Description")

		cmd.Flags().String("name", "", "Name")

		cmd.Flags().String("home-location", "", "Home Location")
		_ = cmd.RegisterFlagCompletionFunc("home-location", cmpl.SuggestCandidatesF(client.Location().Names))

		cmd.Flags().String("server", "", "Server to assign Floating IP to")
		_ = cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(client.Server().Names))

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		_ = cmd.RegisterFlagCompletionFunc("enable-protection", cmpl.SuggestCandidates("delete"))

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, _ []string) (*hcloud.FloatingIP, any, error) {
		typ, _ := cmd.Flags().GetString("type")
		if typ == "" {
			return nil, nil, errors.New("type is required")
		}

		homeLocation, _ := cmd.Flags().GetString("home-location")
		server, _ := cmd.Flags().GetString("server")
		if homeLocation == "" && server == "" {
			return nil, nil, errors.New("one of --home-location or --server is required")
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		serverNameOrID, _ := cmd.Flags().GetString("server")
		labels, _ := cmd.Flags().GetStringToString("label")
		protection, _ := cmd.Flags().GetStringSlice("enable-protection")

		protectionOpts, err := ChangeProtectionCmds.GetChangeProtectionOpts(true, protection)
		if err != nil {
			return nil, nil, err
		}

		createOpts := hcloud.FloatingIPCreateOpts{
			Type:        hcloud.FloatingIPType(typ),
			Description: &description,
			Labels:      labels,
		}
		if name != "" {
			createOpts.Name = &name
		}
		if homeLocation != "" {
			createOpts.HomeLocation = &hcloud.Location{Name: homeLocation}
		}
		if serverNameOrID != "" {
			server, _, err := s.Client().Server().Get(s, serverNameOrID)
			if err != nil {
				return nil, nil, err
			}
			if server == nil {
				return nil, nil, fmt.Errorf("Server not found: %s", serverNameOrID)
			}
			createOpts.Server = server
		}

		result, _, err := s.Client().FloatingIP().Create(s, createOpts)
		if err != nil {
			return nil, nil, err
		}

		if result.Action != nil {
			if err := s.WaitForActions(s, cmd, result.Action); err != nil {
				return nil, nil, err
			}
		}

		cmd.Printf("Floating IP %d created\n", result.FloatingIP.ID)

		if protectionOpts.Delete != nil {
			if err := ChangeProtectionCmds.ChangeProtection(s, cmd, result.FloatingIP, true, protectionOpts); err != nil {
				return nil, nil, err
			}
		}

		floatingIP, _, err := s.Client().FloatingIP().GetByID(s, result.FloatingIP.ID)
		if err != nil {
			return nil, nil, err
		}
		if floatingIP == nil {
			return nil, nil, fmt.Errorf("Floating IP not found: %d", result.FloatingIP.ID)
		}

		return floatingIP, util.Wrap("floating_ip", hcloud.SchemaFromFloatingIP(floatingIP)), nil
	},

	PrintResource: func(_ state.State, cmd *cobra.Command, floatingIP *hcloud.FloatingIP) {
		cmd.Printf("IP%s: %s\n", floatingIP.Type[2:], floatingIP.IP)
	},
}
