package primaryip

import (
	"context"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "Primary IP",
	ShortDescription:     "Describe an Primary IP",
	JSONKeyGetByID:       "primary_ip",
	JSONKeyGetByName:     "primary_ips",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PrimaryIP().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.PrimaryIP().Get(ctx, idOrName)
	},
	PrintText: func(_ context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}) error {
		primaryIP := resource.(*hcloud.PrimaryIP)

		fmt.Printf("ID:\t\t%d\n", primaryIP.ID)
		fmt.Printf("Name:\t\t%s\n", primaryIP.Name)
		fmt.Printf("Created:\t%s (%s)\n", util.Datetime(primaryIP.Created), humanize.Time(primaryIP.Created))
		fmt.Printf("Type:\t\t%s\n", primaryIP.Type)
		fmt.Printf("IP:\t\t%s\n", primaryIP.IP.String())
		fmt.Printf("Blocked:\t%s\n", util.YesNo(primaryIP.Blocked))
		fmt.Printf("Auto delete:\t%s\n", util.YesNo(primaryIP.AutoDelete))
		if primaryIP.AssigneeID != 0 {
			fmt.Printf("Assignee:\n")
			fmt.Printf("  ID:\t%d\n", primaryIP.AssigneeID)
			fmt.Printf("  Type:\t%s\n", primaryIP.AssigneeType)
		} else {
			fmt.Print("Assignee:\n  Not assigned\n")
		}
		fmt.Print("DNS:\n")
		if len(primaryIP.DNSPtr) == 0 {
			fmt.Print("  No reverse DNS entries\n")
		} else {
			for ip, dns := range primaryIP.DNSPtr {
				fmt.Printf("  %s: %s\n", ip, dns)
			}
		}

		fmt.Printf("Protection:\n")
		fmt.Printf("  Delete:\t%s\n", util.YesNo(primaryIP.Protection.Delete))

		fmt.Print("Labels:\n")
		if len(primaryIP.Labels) == 0 {
			fmt.Print("  No labels\n")
		} else {
			for key, value := range primaryIP.Labels {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
		fmt.Printf("Datacenter:\n")
		fmt.Printf("  ID:\t\t%d\n", primaryIP.Datacenter.ID)
		fmt.Printf("  Name:\t\t%s\n", primaryIP.Datacenter.Name)
		fmt.Printf("  Description:\t%s\n", primaryIP.Datacenter.Description)
		fmt.Printf("  Location:\n")
		fmt.Printf("    Name:\t\t%s\n", primaryIP.Datacenter.Location.Name)
		fmt.Printf("    Description:\t%s\n", primaryIP.Datacenter.Location.Description)
		fmt.Printf("    Country:\t\t%s\n", primaryIP.Datacenter.Location.Country)
		fmt.Printf("    City:\t\t%s\n", primaryIP.Datacenter.Location.City)
		fmt.Printf("    Latitude:\t\t%f\n", primaryIP.Datacenter.Location.Latitude)
		fmt.Printf("    Longitude:\t\t%f\n", primaryIP.Datacenter.Location.Longitude)
		return nil
	},
}
