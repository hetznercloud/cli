package floatingip

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "create FLAGS",
			Short:                 "Create a Floating IP",
			Args:                  cobra.NoArgs,
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("type", "", "Type (ipv4 or ipv6) (required)")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("ipv4", "ipv6"))
		cmd.MarkFlagRequired("type")

		cmd.Flags().String("description", "", "Description")

		cmd.Flags().String("name", "", "Name")

		cmd.Flags().String("home-location", "", "Home location")
		cmd.RegisterFlagCompletionFunc("home-location", cmpl.SuggestCandidatesF(client.Location().Names))

		cmd.Flags().String("server", "", "Server to assign Floating IP to")
		cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(client.Server().Names))

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		cmd.RegisterFlagCompletionFunc("enable-protection", cmpl.SuggestCandidates("delete"))

		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		typ, _ := cmd.Flags().GetString("type")
		if typ == "" {
			return errors.New("type is required")
		}

		homeLocation, _ := cmd.Flags().GetString("home-location")
		server, _ := cmd.Flags().GetString("server")
		if homeLocation == "" && server == "" {
			return errors.New("one of --home-location or --server is required")
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		serverNameOrID, _ := cmd.Flags().GetString("server")
		labels, _ := cmd.Flags().GetStringToString("label")
		protection, _ := cmd.Flags().GetStringSlice("enable-protection")

		protectionOps, err := getChangeProtectionOpts(true, protection)
		if err != nil {
			return err
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
			server, _, err := client.Server().Get(ctx, serverNameOrID)
			if err != nil {
				return err
			}
			if server == nil {
				return fmt.Errorf("server not found: %s", serverNameOrID)
			}
			createOpts.Server = server
		}

		result, _, err := client.FloatingIP().Create(ctx, createOpts)
		if err != nil {
			return err
		}

		fmt.Printf("Floating IP %d created\n", result.FloatingIP.ID)

		return changeProtection(ctx, client, waiter, result.FloatingIP, true, protectionOps)
	},
}
