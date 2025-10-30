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
		_, _ = fmt.Fprintf(out, "ID:\t%d\n", server.ID)
		_, _ = fmt.Fprintf(out, "Name:\t%s\n", server.Name)
		_, _ = fmt.Fprintf(out, "Status:\t%s\n", server.Status)
		_, _ = fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(server.Created), humanize.Time(server.Created))

		serverTypeDescription, _ := servertype.DescribeServerType(s, server.ServerType, true)
		_, _ = fmt.Fprintf(out, "Server Type:\t\n")
		_, _ = fmt.Fprintf(out, "%s", util.PrefixLines(serverTypeDescription, "  "))

		// As we already know the location the server is in, we can show the deprecation info
		// of that server type in that specific location.
		locationInfoIndex := slices.IndexFunc(server.ServerType.Locations, func(locInfo hcloud.ServerTypeLocation) bool {
			return locInfo.Location.Name == server.Datacenter.Location.Name
		})
		if locationInfoIndex >= 0 {
			if text := util.DescribeDeprecation(server.ServerType.Locations[locationInfoIndex]); text != "" {
				_, _ = fmt.Fprintf(out, "%s", util.PrefixLines(text, "  "))
			}
		}

		_, _ = fmt.Fprintf(out, "Public Net:\t\n")
		if server.PublicNet.IPv4.IsUnspecified() {
			_, _ = fmt.Fprintf(out, "  IPv4:\tNo Primary IPv4\n")
		} else {
			_, _ = fmt.Fprintf(out, "  IPv4:\t\n")
			_, _ = fmt.Fprintf(out, "    ID:\t%d\n", server.PublicNet.IPv4.ID)
			_, _ = fmt.Fprintf(out, "    IP:\t%s\n", server.PublicNet.IPv4.IP)
			_, _ = fmt.Fprintf(out, "    Blocked:\t%s\n", util.YesNo(server.PublicNet.IPv4.Blocked))
			_, _ = fmt.Fprintf(out, "    DNS:\t%s\n", server.PublicNet.IPv4.DNSPtr)
		}

		if server.PublicNet.IPv6.IsUnspecified() {
			_, _ = fmt.Fprintf(out, "  IPv6:\tNo Primary IPv6\n")
		} else {
			_, _ = fmt.Fprintf(out, "  IPv6:\t\n")
			_, _ = fmt.Fprintf(out, "    ID:\t%d\n", server.PublicNet.IPv6.ID)
			_, _ = fmt.Fprintf(out, "    IP:\t%s\n", server.PublicNet.IPv6.Network.String())
			_, _ = fmt.Fprintf(out, "    Blocked:\t%s\n", util.YesNo(server.PublicNet.IPv6.Blocked))
		}

		if len(server.PublicNet.FloatingIPs) > 0 {
			_, _ = fmt.Fprintf(out, "  Floating IPs:\t\n")
			for _, f := range server.PublicNet.FloatingIPs {
				floatingIP, _, err := s.Client().FloatingIP().GetByID(s, f.ID)
				if err != nil {
					return fmt.Errorf("error fetching Floating IP: %w", err)
				}
				_, _ = fmt.Fprintf(out, "  - ID:\t%d\n", floatingIP.ID)
				_, _ = fmt.Fprintf(out, "    Description:\t%s\n", util.NA(floatingIP.Description))
				_, _ = fmt.Fprintf(out, "    IP:\t%s\n", floatingIP.IP)
			}
		} else {
			_, _ = fmt.Fprintf(out, "  Floating IPs:\tNo Floating IPs\n")
		}

		if len(server.PrivateNet) > 0 {
			_, _ = fmt.Fprintf(out, "Private Net:\t\n")
			for _, n := range server.PrivateNet {
				network, _, err := s.Client().Network().GetByID(s, n.Network.ID)
				if err != nil {
					return fmt.Errorf("error fetching Network: %w", err)
				}
				_, _ = fmt.Fprintf(out, "  - ID:\t%d\n", network.ID)
				_, _ = fmt.Fprintf(out, "    Name:\t%s\n", network.Name)
				_, _ = fmt.Fprintf(out, "    IP:\t%s\n", n.IP.String())
				_, _ = fmt.Fprintf(out, "    MAC Address:\t%s\n", n.MACAddress)
				if len(n.Aliases) > 0 {
					_, _ = fmt.Fprintf(out, "    Alias IPs:\t\n")
					for _, a := range n.Aliases {
						_, _ = fmt.Fprintf(out, "     -\t%s\n", a)
					}
				} else {
					_, _ = fmt.Fprintf(out, "    Alias IPs:\t%s\n", util.NA(""))
				}
			}
		} else {
			_, _ = fmt.Fprintf(out, "Private Net:\tNo Private Networks\n")
		}

		if len(server.Volumes) > 0 {
			_, _ = fmt.Fprintf(out, "Volumes:\t\n")
			for _, v := range server.Volumes {
				volume, _, err := s.Client().Volume().GetByID(s, v.ID)
				if err != nil {
					return fmt.Errorf("error fetching Volume: %w", err)
				}
				_, _ = fmt.Fprintf(out, "  - ID:\t%d\n", volume.ID)
				_, _ = fmt.Fprintf(out, "    Name:\t%s\n", volume.Name)
				_, _ = fmt.Fprintf(out, "    Size:\t%s\n", humanize.Bytes(uint64(volume.Size)*humanize.GByte))
			}
		} else {
			_, _ = fmt.Fprintf(out, "Volumes:\tNo Volumes\n")
		}

		if server.Image != nil {
			_, _ = fmt.Fprintf(out, "Image:\t\n")
			_, _ = fmt.Fprintf(out, "%s", util.PrefixLines(image.DescribeImage(server.Image), "  "))
		} else {
			_, _ = fmt.Fprintf(out, "Image:\tNo Image\n")
		}

		_, _ = fmt.Fprintf(out, "Datacenter:\t\n")
		_, _ = fmt.Fprintf(out, "%s", util.PrefixLines(datacenter.DescribeDatacenter(s.Client(), server.Datacenter, true), "  "))

		_, _ = fmt.Fprintf(out, "Traffic:\t\n")
		_, _ = fmt.Fprintf(out, "  Outgoing:\t%v\n", humanize.IBytes(server.OutgoingTraffic))
		_, _ = fmt.Fprintf(out, "  Ingoing:\t%v\n", humanize.IBytes(server.IngoingTraffic))
		_, _ = fmt.Fprintf(out, "  Included:\t%v\n", humanize.IBytes(server.IncludedTraffic))

		if server.BackupWindow != "" {
			_, _ = fmt.Fprintf(out, "Backup Window:\t%s\n", server.BackupWindow)
		} else {
			_, _ = fmt.Fprintf(out, "Backup Window:\tBackups disabled\n")
		}

		if server.RescueEnabled {
			_, _ = fmt.Fprintf(out, "Rescue System:\tenabled\n")
		} else {
			_, _ = fmt.Fprintf(out, "Rescue System:\tdisabled\n")
		}

		if server.ISO != nil {
			_, _ = fmt.Fprintf(out, "ISO:\t\n")
			_, _ = fmt.Fprintf(out, "%s", util.PrefixLines(iso.DescribeISO(server.ISO), "  "))
		} else {
			_, _ = fmt.Fprintf(out, "ISO:\tNo ISO attached\n")
		}

		_, _ = fmt.Fprintf(out, "Protection:\t\n")
		_, _ = fmt.Fprintf(out, "  Delete:\t%s\n", util.YesNo(server.Protection.Delete))
		_, _ = fmt.Fprintf(out, "  Rebuild:\t%s\n", util.YesNo(server.Protection.Rebuild))

		util.DescribeLabels(out, server.Labels, "")

		if server.PlacementGroup != nil {
			_, _ = fmt.Fprintf(out, "Placement Group:\t\n")
			_, _ = fmt.Fprintf(out, "%s", util.PrefixLines(placementgroup.DescribePlacementGroup(s.Client(), server.PlacementGroup), "  "))
		} else {
			_, _ = fmt.Fprintf(out, "Placement Group:\tNo Placement Group set\n")
		}

		return nil
	},
}
