package floatingip

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

var describeCmd = base.DescribeCmd{
	ResourceNameSingular: "Floating IP",
	ShortDescription:     "Describe an Floating IP",
	JSONKeyGetByID:       "floating_ip",
	JSONKeyGetByName:     "floating_ips",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.FloatingIP().Get(ctx, idOrName)
	},
	PrintText: func(_ context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}) error {
		floatingIP := resource.(*hcloud.FloatingIP)

		fmt.Printf("ID:\t\t%d\n", floatingIP.ID)
		fmt.Printf("Type:\t\t%s\n", floatingIP.Type)
		fmt.Printf("Name:\t\t%s\n", floatingIP.Name)
		fmt.Printf("Description:\t%s\n", util.NA(floatingIP.Description))
		fmt.Printf("Created:\t%s (%s)\n", util.Datetime(floatingIP.Created), humanize.Time(floatingIP.Created))
		if floatingIP.Network != nil {
			fmt.Printf("IP:\t\t%s\n", floatingIP.Network.String())
		} else {
			fmt.Printf("IP:\t\t%s\n", floatingIP.IP.String())
		}
		fmt.Printf("Blocked:\t%s\n", util.YesNo(floatingIP.Blocked))
		fmt.Printf("Home Location:\t%s\n", floatingIP.HomeLocation.Name)
		if floatingIP.Server != nil {
			fmt.Printf("Server:\n")
			fmt.Printf("  ID:\t%d\n", floatingIP.Server.ID)
			fmt.Printf("  Name:\t%s\n", client.Server().ServerName(floatingIP.Server.ID))
		} else {
			fmt.Print("Server:\n  Not assigned\n")
		}
		fmt.Print("DNS:\n")
		if len(floatingIP.DNSPtr) == 0 {
			fmt.Print("  No reverse DNS entries\n")
		} else {
			for ip, dns := range floatingIP.DNSPtr {
				fmt.Printf("  %s: %s\n", ip, dns)
			}
		}

		fmt.Printf("Protection:\n")
		fmt.Printf("  Delete:\t%s\n", util.YesNo(floatingIP.Protection.Delete))

		fmt.Print("Labels:\n")
		if len(floatingIP.Labels) == 0 {
			fmt.Print("  No labels\n")
		} else {
			for key, value := range floatingIP.Labels {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
		return nil
	},
}
