package primaryip

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/datacenter"
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
		_, _ = fmt.Fprintf(out, "ID:\t%d\n", primaryIP.ID)
		_, _ = fmt.Fprintf(out, "Name:\t%s\n", primaryIP.Name)
		_, _ = fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(primaryIP.Created), humanize.Time(primaryIP.Created))
		_, _ = fmt.Fprintf(out, "Type:\t%s\n", primaryIP.Type)
		_, _ = fmt.Fprintf(out, "IP:\t%s\n", primaryIP.IP.String())
		_, _ = fmt.Fprintf(out, "Blocked:\t%s\n", util.YesNo(primaryIP.Blocked))
		_, _ = fmt.Fprintf(out, "Auto delete:\t%s\n", util.YesNo(primaryIP.AutoDelete))

		if primaryIP.AssigneeID != 0 {
			_, _ = fmt.Fprintf(out, "Assignee:\t\n")
			_, _ = fmt.Fprintf(out, "  ID:\t%d\n", primaryIP.AssigneeID)
			_, _ = fmt.Fprintf(out, "  Type:\t%s\n", primaryIP.AssigneeType)
		} else {
			_, _ = fmt.Fprintf(out, "Assignee:\tNot assigned\n")
		}

		if len(primaryIP.DNSPtr) == 0 {
			_, _ = fmt.Fprintf(out, "DNS:\tNo reverse DNS entries\n")
		} else {
			_, _ = fmt.Fprintf(out, "DNS:\t\n")
			for ip, dns := range primaryIP.DNSPtr {
				_, _ = fmt.Fprintf(out, "  %s:\t%s\n", ip, dns)
			}
		}

		_, _ = fmt.Fprintf(out, "Protection:\t\n")
		_, _ = fmt.Fprintf(out, "  Delete:\t%s\n", util.YesNo(primaryIP.Protection.Delete))

		util.DescribeLabels(out, primaryIP.Labels, "")

		_, _ = fmt.Fprintf(out, "Datacenter:\t\n")
		_, _ = fmt.Fprintf(out, "%s", util.PrefixLines(datacenter.DescribeDatacenter(s.Client(), primaryIP.Datacenter, true), "  "))
		return nil
	},
}
