package primaryip

import (
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.PrimaryIP]{
	ResourceNameSingular: "Primary IP",
	ShortDescription:     "Describe an Primary IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PrimaryIP().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.PrimaryIP, any, error) {
		ip, _, err := s.Client().PrimaryIP().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return ip, hcloud.SchemaFromPrimaryIP(ip), nil
	},
	PrintText: func(_ state.State, cmd *cobra.Command, primaryIP *hcloud.PrimaryIP) error {
		cmd.Printf("ID:\t\t%d\n", primaryIP.ID)
		cmd.Printf("Name:\t\t%s\n", primaryIP.Name)
		cmd.Printf("Created:\t%s (%s)\n", util.Datetime(primaryIP.Created), humanize.Time(primaryIP.Created))
		cmd.Printf("Type:\t\t%s\n", primaryIP.Type)
		cmd.Printf("IP:\t\t%s\n", primaryIP.IP.String())
		cmd.Printf("Blocked:\t%s\n", util.YesNo(primaryIP.Blocked))
		cmd.Printf("Auto delete:\t%s\n", util.YesNo(primaryIP.AutoDelete))
		if primaryIP.AssigneeID != 0 {
			cmd.Printf("Assignee:\n")
			cmd.Printf("  ID:\t%d\n", primaryIP.AssigneeID)
			cmd.Printf("  Type:\t%s\n", primaryIP.AssigneeType)
		} else {
			cmd.Print("Assignee:\n  Not assigned\n")
		}
		cmd.Print("DNS:\n")
		if len(primaryIP.DNSPtr) == 0 {
			cmd.Print("  No reverse DNS entries\n")
		} else {
			for ip, dns := range primaryIP.DNSPtr {
				cmd.Printf("  %s: %s\n", ip, dns)
			}
		}

		cmd.Printf("Protection:\n")
		cmd.Printf("  Delete:\t%s\n", util.YesNo(primaryIP.Protection.Delete))

		cmd.Print("Labels:\n")
		if len(primaryIP.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range util.IterateInOrder(primaryIP.Labels) {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}
		cmd.Printf("Datacenter:\n")
		cmd.Printf("  ID:\t\t%d\n", primaryIP.Datacenter.ID)
		cmd.Printf("  Name:\t\t%s\n", primaryIP.Datacenter.Name)
		cmd.Printf("  Description:\t%s\n", primaryIP.Datacenter.Description)
		cmd.Printf("  Location:\n")
		cmd.Printf("    Name:\t\t%s\n", primaryIP.Datacenter.Location.Name)
		cmd.Printf("    Description:\t%s\n", primaryIP.Datacenter.Location.Description)
		cmd.Printf("    Country:\t\t%s\n", primaryIP.Datacenter.Location.Country)
		cmd.Printf("    City:\t\t%s\n", primaryIP.Datacenter.Location.City)
		cmd.Printf("    Latitude:\t\t%f\n", primaryIP.Datacenter.Location.Latitude)
		cmd.Printf("    Longitude:\t\t%f\n", primaryIP.Datacenter.Location.Longitude)
		return nil
	},
}
