package primaryip

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd{
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
	Run: func(s state.State, cmd *cobra.Command, args []string) (any, any, error) {
		typ, _ := cmd.Flags().GetString("type")
		name, _ := cmd.Flags().GetString("name")
		assigneeID, _ := cmd.Flags().GetInt64("assignee-id")
		datacenter, _ := cmd.Flags().GetString("datacenter")
		protection, _ := cmd.Flags().GetStringSlice("enable-protection")

		protectionOpts, err := getChangeProtectionOpts(true, protection)
		if err != nil {
			return nil, nil, err
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

		result, _, err := s.PrimaryIP().Create(s, createOpts)
		if err != nil {
			return nil, nil, err
		}

		if result.Action != nil {
			if err := s.ActionProgress(cmd, s, result.Action); err != nil {
				return nil, nil, err
			}
		}

		cmd.Printf("Primary IP %d created\n", result.PrimaryIP.ID)

		if len(protection) > 0 {
			if err := changeProtection(s, cmd, result.PrimaryIP, true, protectionOpts); err != nil {
				return nil, nil, err
			}
		}

		return result.PrimaryIP, util.Wrap("primary_ip", hcloud.SchemaFromPrimaryIP(result.PrimaryIP)), nil
	},
	PrintResource: func(_ state.State, cmd *cobra.Command, resource any) {
		primaryIP := resource.(*hcloud.PrimaryIP)
		cmd.Printf("IP%s: %s\n", primaryIP.Type[2:], primaryIP.IP)
	},
}
