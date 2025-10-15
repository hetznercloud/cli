package primaryip

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd[*hcloud.PrimaryIP]{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "create [options] --type <ipv4|ipv6> --name <name>",
			Short:                 "Create a Primary IP",
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("type", "", "Type (ipv4 or ipv6) (required)")
		_ = cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("ipv4", "ipv6"))
		_ = cmd.MarkFlagRequired("type")

		cmd.Flags().String("name", "", "Name (required)")
		_ = cmd.MarkFlagRequired("name")

		cmd.Flags().Int64("assignee-id", 0, "Assignee (usually a Server) to assign Primary IP to (required if 'datacenter' is not specified)")

		cmd.Flags().String("datacenter", "", "Datacenter (ID or name) (required if 'assignee-id' is not specified)")
		_ = cmd.RegisterFlagCompletionFunc("datacenter", cmpl.SuggestCandidatesF(client.Datacenter().Names))

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		_ = cmd.RegisterFlagCompletionFunc("enable-protection", cmpl.SuggestCandidates("delete"))

		cmd.Flags().Bool("auto-delete", false, "Delete Primary IP if assigned resource is deleted (true, false)")

		cmd.MarkFlagsOneRequired("assignee-id", "datacenter")
		cmd.MarkFlagsMutuallyExclusive("assignee-id", "datacenter")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, _ []string) (*hcloud.PrimaryIP, any, error) {
		typ, _ := cmd.Flags().GetString("type")
		name, _ := cmd.Flags().GetString("name")
		assigneeID, _ := cmd.Flags().GetInt64("assignee-id")
		datacenter, _ := cmd.Flags().GetString("datacenter")
		labels, _ := cmd.Flags().GetStringToString("label")
		protection, _ := cmd.Flags().GetStringSlice("enable-protection")
		autoDelete, _ := cmd.Flags().GetBool("auto-delete")

		protectionOpts, err := getChangeProtectionOpts(true, protection)
		if err != nil {
			return nil, nil, err
		}

		createOpts := hcloud.PrimaryIPCreateOpts{
			Type:         hcloud.PrimaryIPType(typ),
			Name:         name,
			AssigneeType: "server",
			Datacenter:   datacenter,
			Labels:       labels,
		}
		if assigneeID != 0 {
			createOpts.AssigneeID = &assigneeID
		}
		if cmd.Flags().Changed("auto-delete") {
			createOpts.AutoDelete = &autoDelete
		}

		result, _, err := s.Client().PrimaryIP().Create(s, createOpts)
		if err != nil {
			return nil, nil, err
		}

		if result.Action != nil {
			if err := s.WaitForActions(s, cmd, result.Action); err != nil {
				return nil, nil, err
			}
		}

		cmd.Printf("Primary IP %d created\n", result.PrimaryIP.ID)

		if len(protection) > 0 {
			if err := changeProtection(s, cmd, result.PrimaryIP, true, protectionOpts); err != nil {
				return nil, nil, err
			}
		}

		primaryIP, _, err := s.Client().PrimaryIP().GetByID(s, result.PrimaryIP.ID)
		if err != nil {
			return nil, nil, err
		}
		if primaryIP == nil {
			return nil, nil, fmt.Errorf("Primary IP not found: %d", result.PrimaryIP.ID)
		}

		return primaryIP, util.Wrap("primary_ip", hcloud.SchemaFromPrimaryIP(primaryIP)), nil
	},
	PrintResource: func(_ state.State, cmd *cobra.Command, primaryIP *hcloud.PrimaryIP) {
		cmd.Printf("IP%s: %s\n", primaryIP.Type[2:], primaryIP.IP)
	},
}
