package floatingip

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.FloatingIP]{
	ResourceNameSingular: "Floating IP",
	ShortDescription:     "Describe a Floating IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.FloatingIP, any, error) {
		ip, _, err := s.Client().FloatingIP().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return ip, hcloud.SchemaFromFloatingIP(ip), nil
	},
	PrintText: func(s state.State, _ *cobra.Command, out io.Writer, floatingIP *hcloud.FloatingIP) error {
		fmt.Fprintf(out, "ID:\t%d\n", floatingIP.ID)
		fmt.Fprintf(out, "Type:\t%s\n", floatingIP.Type)
		fmt.Fprintf(out, "Name:\t%s\n", floatingIP.Name)
		fmt.Fprintf(out, "Description:\t%s\n", util.NA(floatingIP.Description))
		fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(floatingIP.Created), humanize.Time(floatingIP.Created))

		if floatingIP.Network != nil {
			fmt.Fprintf(out, "IP:\t%s\n", floatingIP.Network.String())
		} else {
			fmt.Fprintf(out, "IP:\t%s\n", floatingIP.IP.String())
		}

		fmt.Fprintf(out, "Blocked:\t%s\n", util.YesNo(floatingIP.Blocked))
		fmt.Fprintf(out, "Home Location:\t%s\n", floatingIP.HomeLocation.Name)

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Server:\n")
		if floatingIP.Server != nil {
			fmt.Fprintf(out, "  ID:\t%d\n", floatingIP.Server.ID)
			fmt.Fprintf(out, "  Name:\t%s\n", s.Client().Server().ServerName(floatingIP.Server.ID))
		} else {
			fmt.Fprintf(out, "  Not assigned\n")
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "DNS:\n")
		if len(floatingIP.DNSPtr) == 0 {
			fmt.Fprintf(out, "  No reverse DNS entries\n")
		} else {
			for ip, dns := range floatingIP.DNSPtr {
				fmt.Fprintf(out, "  %s:\t%s\n", ip, dns)
			}
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Protection:\n")
		fmt.Fprintf(out, "  Delete:\t%s\n", util.YesNo(floatingIP.Protection.Delete))

		fmt.Fprintln(out)
		util.DescribeLabels(out, floatingIP.Labels, "")
		return nil
	},
}
