package server

import (
	"fmt"
	"io"
	"slices"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/datacenter"
	"github.com/hetznercloud/cli/internal/cmd/image"
	"github.com/hetznercloud/cli/internal/cmd/iso"
	"github.com/hetznercloud/cli/internal/cmd/placementgroup"
	"github.com/hetznercloud/cli/internal/cmd/servertype"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.Server]{
	ResourceNameSingular: "Server",
	ShortDescription:     "Describe a Server",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Server().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Server, any, error) {
		srv, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return srv, hcloud.SchemaFromServer(srv), nil
	},
	PrintText: func(s state.State, _ *cobra.Command, out io.Writer, server *hcloud.Server) error {
		fmt.Fprintf(out, "ID:\t%d\n", server.ID)
		fmt.Fprintf(out, "Name:\t%s\n", server.Name)
		fmt.Fprintf(out, "Status:\t%s\n", server.Status)
		fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(server.Created), humanize.Time(server.Created))

		serverTypeDescription, _ := servertype.DescribeServerType(s, server.ServerType, true)
		fmt.Fprintln(out)
		fmt.Fprintf(out, "Server Type:\n")
		fmt.Fprintf(out, "%s", util.PrefixLines(serverTypeDescription, "  "))

		// As we already know the location the server is in, we can show the deprecation info
		// of that server type in that specific location.
		locationInfoIndex := slices.IndexFunc(server.ServerType.Locations, func(locInfo hcloud.ServerTypeLocation) bool {
			return locInfo.Location.Name == server.Datacenter.Location.Name
		})
		if locationInfoIndex >= 0 {
			if text := util.DescribeDeprecation(server.ServerType.Locations[locationInfoIndex]); text != "" {
				fmt.Fprintln(out)
				fmt.Fprint(out, util.PrefixLines(text, "  "))
			}
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Public Net:\n")
		fmt.Fprintf(out, "  IPv4:\n")
		if server.PublicNet.IPv4.IsUnspecified() {
			fmt.Fprintf(out, "    No Primary IPv4\n")
		} else {
			fmt.Fprintf(out, "    ID:\t%d\n", server.PublicNet.IPv4.ID)
			fmt.Fprintf(out, "    IP:\t%s\n", server.PublicNet.IPv4.IP)
			fmt.Fprintf(out, "    Blocked:\t%s\n", util.YesNo(server.PublicNet.IPv4.Blocked))
			fmt.Fprintf(out, "    DNS:\t%s\n", server.PublicNet.IPv4.DNSPtr)
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "  IPv6:\n")
		if server.PublicNet.IPv6.IsUnspecified() {
			fmt.Fprintf(out, "    No Primary IPv6\n")
		} else {
			fmt.Fprintf(out, "    ID:\t%d\n", server.PublicNet.IPv6.ID)
			fmt.Fprintf(out, "    IP:\t%s\n", server.PublicNet.IPv6.Network.String())
			fmt.Fprintf(out, "    Blocked:\t%s\n", util.YesNo(server.PublicNet.IPv6.Blocked))
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "  Floating IPs:\n")
		if len(server.PublicNet.FloatingIPs) > 0 {
			for _, f := range server.PublicNet.FloatingIPs {
				floatingIP, _, err := s.Client().FloatingIP().GetByID(s, f.ID)
				if err != nil {
					return fmt.Errorf("error fetching Floating IP: %w", err)
				}
				fmt.Fprintf(out, "  - ID:\t%d\n", floatingIP.ID)
				fmt.Fprintf(out, "    Description:\t%s\n", util.NA(floatingIP.Description))
				fmt.Fprintf(out, "    IP:\t%s\n", floatingIP.IP)
			}
		} else {
			fmt.Fprintf(out, "    No Floating IPs\n")
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Private Net:\n")
		if len(server.PrivateNet) > 0 {
			for i, n := range server.PrivateNet {
				if i > 0 {
					fmt.Fprintln(out)
				}
				network, _, err := s.Client().Network().GetByID(s, n.Network.ID)
				if err != nil {
					return fmt.Errorf("error fetching Network: %w", err)
				}
				fmt.Fprintf(out, "  - ID:\t%d\n", network.ID)
				fmt.Fprintf(out, "    Name:\t%s\n", network.Name)
				fmt.Fprintf(out, "    IP:\t%s\n", n.IP.String())
				fmt.Fprintf(out, "    MAC Address:\t%s\n", n.MACAddress)
				if len(n.Aliases) > 0 {
					fmt.Fprintf(out, "    Alias IPs:\t\n")
					for _, a := range n.Aliases {
						fmt.Fprintf(out, "     -\t%s\n", a)
					}
				} else {
					fmt.Fprintf(out, "    Alias IPs:\t%s\n", util.NA(""))
				}
			}
		} else {
			fmt.Fprintf(out, "  No Private Networks\n")
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Volumes:\n")
		if len(server.Volumes) > 0 {
			for _, v := range server.Volumes {
				volume, _, err := s.Client().Volume().GetByID(s, v.ID)
				if err != nil {
					return fmt.Errorf("error fetching Volume: %w", err)
				}
				fmt.Fprintf(out, "  - ID:\t%d\n", volume.ID)
				fmt.Fprintf(out, "    Name:\t%s\n", volume.Name)
				fmt.Fprintf(out, "    Size:\t%s\n", humanize.Bytes(uint64(volume.Size)*humanize.GByte))
			}
		} else {
			fmt.Fprintf(out, "  No Volumes\n")
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Image:\n")
		if server.Image != nil {
			fmt.Fprint(out, util.PrefixLines(image.DescribeImage(server.Image), "  "))
		} else {
			fmt.Fprintf(out, "  No Image\n")
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Datacenter:\n")
		fmt.Fprint(out, util.PrefixLines(datacenter.DescribeDatacenter(s.Client(), server.Datacenter, true), "  "))

		fmt.Fprintln(out)
		if server.BackupWindow != "" {
			fmt.Fprintf(out, "Backup Window:\t%s\n", server.BackupWindow)
		} else {
			fmt.Fprintf(out, "Backup Window:\tBackups disabled\n")
		}

		fmt.Fprintln(out)
		if server.RescueEnabled {
			fmt.Fprintf(out, "Rescue System:\tenabled\n")
		} else {
			fmt.Fprintf(out, "Rescue System:\tdisabled\n")
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Traffic:\n")
		fmt.Fprintf(out, "  Outgoing:\t%v\n", humanize.IBytes(server.OutgoingTraffic))
		fmt.Fprintf(out, "  Ingoing:\t%v\n", humanize.IBytes(server.IngoingTraffic))
		fmt.Fprintf(out, "  Included:\t%v\n", humanize.IBytes(server.IncludedTraffic))

		fmt.Fprintln(out)
		fmt.Fprintf(out, "ISO:\n")
		if server.ISO != nil {
			fmt.Fprint(out, util.PrefixLines(iso.DescribeISO(server.ISO), "  "))
		} else {
			fmt.Fprintf(out, "  No ISO attached\n")
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Protection:\n")
		fmt.Fprintf(out, "  Delete:\t%s\n", util.YesNo(server.Protection.Delete))
		fmt.Fprintf(out, "  Rebuild:\t%s\n", util.YesNo(server.Protection.Rebuild))

		fmt.Fprintln(out)
		util.DescribeLabels(out, server.Labels, "")

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Placement Group:\n")
		if server.PlacementGroup != nil {
			fmt.Fprint(out, util.PrefixLines(placementgroup.DescribePlacementGroup(s.Client(), server.PlacementGroup), "  "))
		} else {
			fmt.Fprintf(out, "  No Placement Group set\n")
		}

		return nil
	},
}
