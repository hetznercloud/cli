package floatingip

import (
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "Floating IP",
	ShortDescription:     "Describe an Floating IP",
	JSONKeyGetByID:       "floating_ip",
	JSONKeyGetByName:     "floating_ips",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, interface{}, error) {
		ip, _, err := s.Client().FloatingIP().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return ip, hcloud.SchemaFromFloatingIP(ip), nil
	},
	PrintText: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		floatingIP := resource.(*hcloud.FloatingIP)

		cmd.Printf("ID:\t\t%d\n", floatingIP.ID)
		cmd.Printf("Type:\t\t%s\n", floatingIP.Type)
		cmd.Printf("Name:\t\t%s\n", floatingIP.Name)
		cmd.Printf("Description:\t%s\n", util.NA(floatingIP.Description))
		cmd.Printf("Created:\t%s (%s)\n", util.Datetime(floatingIP.Created), humanize.Time(floatingIP.Created))
		if floatingIP.Network != nil {
			cmd.Printf("IP:\t\t%s\n", floatingIP.Network.String())
		} else {
			cmd.Printf("IP:\t\t%s\n", floatingIP.IP.String())
		}
		cmd.Printf("Blocked:\t%s\n", util.YesNo(floatingIP.Blocked))
		cmd.Printf("Home Location:\t%s\n", floatingIP.HomeLocation.Name)
		if floatingIP.Server != nil {
			cmd.Printf("Server:\n")
			cmd.Printf("  ID:\t%d\n", floatingIP.Server.ID)
			cmd.Printf("  Name:\t%s\n", s.Client().Server().ServerName(floatingIP.Server.ID))
		} else {
			cmd.Print("Server:\n  Not assigned\n")
		}
		cmd.Print("DNS:\n")
		if len(floatingIP.DNSPtr) == 0 {
			cmd.Print("  No reverse DNS entries\n")
		} else {
			for ip, dns := range floatingIP.DNSPtr {
				cmd.Printf("  %s: %s\n", ip, dns)
			}
		}

		cmd.Printf("Protection:\n")
		cmd.Printf("  Delete:\t%s\n", util.YesNo(floatingIP.Protection.Delete))

		cmd.Print("Labels:\n")
		if len(floatingIP.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range floatingIP.Labels {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}
		return nil
	},
}
