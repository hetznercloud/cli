package network

import (
	"fmt"
	"net"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var networkMemberTypeStrings = []string{
	string(hcloud.NetworkMemberTypeLoadBalancer),
	string(hcloud.NetworkMemberTypeServer),
}

var networkMemberStatusStrings = []string{
	string(hcloud.NetworkMemberStatusAttaching),
	string(hcloud.NetworkMemberStatusDetaching),
	string(hcloud.NetworkMemberStatusError),
	string(hcloud.NetworkMemberStatusOK),
	string(hcloud.NetworkMemberStatusUpdating),
}

var ListMembersCmd = &base.ListCmd[*hcloud.NetworkMember, schema.NetworkMember]{
	ResourceNamePlural: "Network Members",
	JSONKeyGetByName:   "members",

	DefaultColumns: []string{"type", "id", "ip", "status", "alias_ips", "subnet"},
	SortOption:     config.OptionSortNetworkMember,

	ValidArgsFunction: func(client hcapi2.Client) cobra.CompletionFunc {
		return cmpl.SuggestCandidatesF(client.Network().Names)
	},

	PositionalArgumentOverride: []string{"network"},

	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().StringSlice("type", nil, "Only Members with one of these types are displayed")
		_ = cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates(networkMemberTypeStrings...))

		cmd.Flags().StringSlice("status", nil, "Only Members with one of these statuses are displayed")
		_ = cmd.RegisterFlagCompletionFunc("status", cmpl.SuggestCandidates(networkMemberStatusStrings...))

		cmd.Flags().StringSlice("subnet", nil, "Only Members attached to one of these subnets are displayed")
	},

	FetchWithArgs: func(s state.State, flags *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string, args []string) ([]*hcloud.NetworkMember, error) {
		networkIDOrName := args[0]

		types, _ := flags.GetStringSlice("type")
		subnets, _ := flags.GetStringSlice("subnet")
		statuses, _ := flags.GetStringSlice("status")

		opts := hcloud.NetworkMemberListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}

		if len(types) > 0 {
			for _, typ := range types {
				if !slices.Contains(networkMemberTypeStrings, typ) {
					return nil, fmt.Errorf("invalid type: %s", typ)
				}
				opts.Type = append(opts.Type, hcloud.NetworkMemberType(typ))
			}
		}

		if len(subnets) > 0 {
			for _, subnet := range subnets {
				_, ipNet, err := net.ParseCIDR(subnet)
				if err != nil {
					return nil, fmt.Errorf("invalid subnet: %s", subnet)
				}
				opts.Subnet = append(opts.Subnet, ipNet)
			}
		}

		if len(statuses) > 0 {
			for _, status := range statuses {
				if !slices.Contains(networkMemberStatusStrings, status) {
					return nil, fmt.Errorf("invalid status: %s", status)
				}
				opts.Status = append(opts.Status, hcloud.NetworkMemberStatus(status))
			}
		}

		network, _, err := s.Client().Network().Get(s, networkIDOrName)
		if err != nil {
			return nil, err
		}
		if network == nil {
			return nil, fmt.Errorf("Network not found: %s", networkIDOrName)
		}

		return s.Client().Network().AllMembersWithOpts(s, network, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.NetworkMember], _ hcapi2.Client) {
		t.
			AddAllowedFields(&hcloud.NetworkMember{}).
			AddFieldFn("ip", func(member *hcloud.NetworkMember) string {
				if member.IP == nil {
					return util.NA("")
				}
				return member.IP.String()
			}).
			AddFieldFn("alias_ips", func(member *hcloud.NetworkMember) string {
				var aliasIPs []string
				for _, aliasIP := range member.AliasIPs {
					aliasIPs = append(aliasIPs, aliasIP.String())
				}
				return util.NA(strings.Join(aliasIPs, "\n"))
			}).
			AddFieldFn("subnet", func(member *hcloud.NetworkMember) string {
				if member.Subnet == nil {
					return util.NA("")
				}
				return member.Subnet.String()
			})
	},

	Schema: hcloud.SchemaFromNetworkMember,

	Configure: func(_ state.State, cmd *cobra.Command) *cobra.Command {
		cmd.Use = "list-members [options] <network>"
		return cmd
	},
}
