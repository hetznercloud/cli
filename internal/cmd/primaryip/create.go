package primaryip

import (
	"context"
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
			Short:                 "Create a Primary IP",
			Args:                  cobra.NoArgs,
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("type", "", "Type (ipv4 or ipv6) (required)")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("ipv4", "ipv6"))
		cmd.MarkFlagRequired("type")

		cmd.Flags().String("name", "", "Name (required)")
		cmd.MarkFlagRequired("name")

		cmd.Flags().Int64("assignee-id", 0, "Assignee (usually a server) to assign Primary IP to")

		cmd.Flags().String("datacenter", "", "Datacenter (ID or name)")
		cmd.RegisterFlagCompletionFunc("datacenter", cmpl.SuggestCandidatesF(client.Datacenter().Names))

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		cmd.RegisterFlagCompletionFunc("enable-protection", cmpl.SuggestCandidates("delete"))

		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		typ, _ := cmd.Flags().GetString("type")
		name, _ := cmd.Flags().GetString("name")
		assigneeID, _ := cmd.Flags().GetInt64("assignee-id")
		datacenter, _ := cmd.Flags().GetString("datacenter")
		protection, _ := cmd.Flags().GetStringSlice("enable-protection")

		protectionOpts, err := getChangeProtectionOpts(true, protection)
		if err != nil {
			return err
		}

		createOpts := hcloud.PrimaryIPCreateOpts{
			Type:         hcloud.PrimaryIPType(typ),
			Name:         name,
			AssigneeType: "server",
			Datacenter:   datacenter,
		}
		if assigneeID != 0 {
			createOpts.AssigneeID = &assigneeID
		}

		result, _, err := client.PrimaryIP().Create(ctx, createOpts)
		if err != nil {
			return err
		}

		if result.Action != nil {
			if err := waiter.ActionProgress(ctx, result.Action); err != nil {
				return err
			}
		}

		fmt.Printf("Primary IP %d created\n", result.PrimaryIP.ID)

		if len(protection) > 0 {
			if err := changeProtection(ctx, client, waiter, result.PrimaryIP, true, protectionOpts); err != nil {
				return err
			}
		}

		return nil
	},
}
