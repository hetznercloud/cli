package primaryip

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/datacenter"
	"github.com/hetznercloud/cli/internal/cmd/location"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.PrimaryIP]{
	ResourceNameSingular: "Primary IP",
	ShortDescription:     "Describe a Primary IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PrimaryIP().Names(false, false, nil) },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.PrimaryIP, any, error) {
		ip, _, err := s.Client().PrimaryIP().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return ip, hcloud.SchemaFromPrimaryIP(ip), nil
	},
	PrintText: func(s state.State, _ *cobra.Command, out io.Writer, primaryIP *hcloud.PrimaryIP) error {
		fmt.Fprintf(out, "ID:\t%d\n", primaryIP.ID)
		fmt.Fprintf(out, "Name:\t%s\n", primaryIP.Name)
		fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(primaryIP.Created), humanize.Time(primaryIP.Created))
		fmt.Fprintf(out, "Type:\t%s\n", primaryIP.Type)
		fmt.Fprintf(out, "IP:\t%s\n", primaryIP.IP.String())
		fmt.Fprintf(out, "Blocked:\t%s\n", util.YesNo(primaryIP.Blocked))
		fmt.Fprintf(out, "Auto delete:\t%s\n", util.YesNo(primaryIP.AutoDelete))

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Assignee:\n")
		if primaryIP.AssigneeID != 0 {
			fmt.Fprintf(out, "  ID:\t%d\n", primaryIP.AssigneeID)
			fmt.Fprintf(out, "  Type:\t%s\n", primaryIP.AssigneeType)
		} else {
			fmt.Fprintf(out, "  Not assigned\n")
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "DNS:\n")
		if len(primaryIP.DNSPtr) == 0 {
			fmt.Fprintf(out, "  No reverse DNS entries\n")
		} else {
			for ip, dns := range primaryIP.DNSPtr {
				fmt.Fprintf(out, "  %s:\t%s\n", ip, dns)
			}
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Protection:\n")
		fmt.Fprintf(out, "  Delete:\t%s\n", util.YesNo(primaryIP.Protection.Delete))

		fmt.Fprintln(out)
		util.DescribeLabels(out, primaryIP.Labels, "")

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Location:\n")
		fmt.Fprintf(out, "%s", util.PrefixLines(location.DescribeLocation(primaryIP.Location), "  "))

		if primaryIP.Datacenter != nil {
			fmt.Fprintln(out)
			fmt.Fprintf(out, "Datacenter:\n")
			fmt.Fprintf(out, "%s", util.PrefixLines(datacenter.DescribeDatacenter(s.Client(), primaryIP.Datacenter, true), "  "))
		}

		return nil
	},
}
