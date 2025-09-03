package server

import (
	"fmt"
	"strconv"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
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
	PrintText: func(s state.State, cmd *cobra.Command, server *hcloud.Server, w base.DescribeWriter) error {

		w.WriteLine("ID:", strconv.FormatInt(server.ID, 10))
		w.WriteLine("Name:", server.Name)
		w.WriteLine("Status:", string(server.Status))
		w.WriteLine("Created:", fmt.Sprintf("%s (%s)", util.Datetime(server.Created), humanize.Time(server.Created)))

		w.WriteLine("Server Type:", fmt.Sprintf("%s (ID: %d)", server.ServerType.Name, server.ServerType.ID))
		servertype.DescribeServerType(server.ServerType, w.NewSubWriter("  "))

		if text := util.DescribeDeprecation(server.ServerType); text != "" {
			cmd.Print(util.PrefixLines(text, "  "))
		}

		w.WriteLine("Public Net:")
		w.WriteLine("  IPv4:")
		if server.PublicNet.IPv4.IsUnspecified() {
			w.WriteLine("    No Primary IPv4")
		} else {
			w.WriteLine("    ID:", strconv.FormatInt(server.PublicNet.IPv4.ID, 10))
			w.WriteLine("    IP:", server.PublicNet.IPv4.IP.String())
			w.WriteLine("    Blocked:", util.YesNo(server.PublicNet.IPv4.Blocked))
			w.WriteLine("    DNS:", server.PublicNet.IPv4.DNSPtr)
		}

		w.WriteLine("  IPv6:")
		if server.PublicNet.IPv6.IsUnspecified() {
			w.WriteLine("    No Primary IPv6")
		} else {
			w.WriteLine("    ID:", strconv.FormatInt(server.PublicNet.IPv6.ID, 10))
			w.WriteLine("    IP:", server.PublicNet.IPv6.Network.String())
			w.WriteLine("    Blocked:", util.YesNo(server.PublicNet.IPv6.Blocked))
		}

		w.WriteLine("  Floating IPs:")
		if len(server.PublicNet.FloatingIPs) > 0 {
			for _, f := range server.PublicNet.FloatingIPs {
				floatingIP, _, err := s.Client().FloatingIP().GetByID(s, f.ID)
				if err != nil {
					return fmt.Errorf("error fetching Floating IP: %w", err)
				}
				w.WriteLine("    - ID:", strconv.FormatInt(floatingIP.ID, 10))
				w.WriteLine("      Description:", util.NA(floatingIP.Description))
				w.WriteLine("      IP:", floatingIP.IP.String())
			}
		} else {
			w.WriteLine("    No Floating IPs")
		}

		w.WriteLine("Private Net:")
		if len(server.PrivateNet) > 0 {
			for _, n := range server.PrivateNet {
				network, _, err := s.Client().Network().GetByID(s, n.Network.ID)
				if err != nil {
					return fmt.Errorf("error fetching Network: %w", err)
				}
				w.WriteLine("  - ID:", strconv.FormatInt(network.ID, 10))
				w.WriteLine("    Name:", network.Name)
				w.WriteLine("    IP:", n.IP.String())
				w.WriteLine("    MAC Address:", n.MACAddress)
				if len(n.Aliases) > 0 {
					w.WriteLine("    Alias IPs:")
					for _, a := range n.Aliases {
						w.WriteLine("     -", a.String())
					}
				} else {
					w.WriteLine("    Alias IPs:", "-")
				}
			}
		} else {
			w.WriteLine("    No Private Networks")
		}

		w.WriteLine("Volumes:")
		if len(server.Volumes) > 0 {
			for _, v := range server.Volumes {
				volume, _, err := s.Client().Volume().GetByID(s, v.ID)
				if err != nil {
					return fmt.Errorf("error fetching Volume: %w", err)
				}
				w.WriteLine("  - ID:", strconv.FormatInt(volume.ID, 10))
				w.WriteLine("    Name:", volume.Name)
				w.WriteLine("    Size:", fmt.Sprintf("%s", humanize.Bytes(uint64(volume.Size)*humanize.GByte)))
			}
		} else {
			w.WriteLine("  No Volumes")
		}

		/*
			w.WriteLine("Image:")
			if server.Image != nil {
				image := server.Image
				cmd.Printf("  ID:\t\t%d\n", image.ID)
				cmd.Printf("  Type:\t\t%s\n", image.Type)
				cmd.Printf("  Status:\t%s\n", image.Status)
				cmd.Printf("  Name:\t\t%s\n", util.NA(image.Name))
				cmd.Printf("  Description:\t%s\n", image.Description)
				if image.ImageSize != 0 {
					cmd.Printf("  Image size:\t%.2f GB\n", image.ImageSize)
				} else {
					cmd.Printf("  Image size:\t%s\n", util.NA(""))
				}
				cmd.Printf("  Disk size:\t%.0f GB\n", image.DiskSize)
				cmd.Printf("  Created:\t%s (%s)\n", util.Datetime(image.Created), humanize.Time(image.Created))
				cmd.Printf("  OS flavor:\t%s\n", image.OSFlavor)
				cmd.Printf("  OS version:\t%s\n", util.NA(image.OSVersion))
				cmd.Printf("  Rapid deploy:\t%s\n", util.YesNo(image.RapidDeploy))
			} else {
				cmd.Printf("  No Image\n")
			}

			cmd.Printf("Datacenter:\n")
			cmd.Printf("  ID:\t\t%d\n", server.Datacenter.ID)
			cmd.Printf("  Name:\t\t%s\n", server.Datacenter.Name)
			cmd.Printf("  Description:\t%s\n", server.Datacenter.Description)
			cmd.Printf("  Location:\n")
			cmd.Printf("    Name:\t\t%s\n", server.Datacenter.Location.Name)
			cmd.Printf("    Description:\t%s\n", server.Datacenter.Location.Description)
			cmd.Printf("    Country:\t\t%s\n", server.Datacenter.Location.Country)
			cmd.Printf("    City:\t\t%s\n", server.Datacenter.Location.City)
			cmd.Printf("    Latitude:\t\t%f\n", server.Datacenter.Location.Latitude)
			cmd.Printf("    Longitude:\t\t%f\n", server.Datacenter.Location.Longitude)

			cmd.Printf("Traffic:\n")
			cmd.Printf("  Outgoing:\t%v\n", humanize.IBytes(server.OutgoingTraffic))
			cmd.Printf("  Ingoing:\t%v\n", humanize.IBytes(server.IngoingTraffic))
			cmd.Printf("  Included:\t%v\n", humanize.IBytes(server.IncludedTraffic))

			if server.BackupWindow != "" {
				cmd.Printf("Backup Window:\t%s\n", server.BackupWindow)
			} else {
				cmd.Printf("Backup Window:\tBackups disabled\n")
			}

			if server.RescueEnabled {
				cmd.Printf("Rescue System:\tenabled\n")
			} else {
				cmd.Printf("Rescue System:\tdisabled\n")
			}

			cmd.Printf("ISO:\n")
			if server.ISO != nil {
				cmd.Printf("  ID:\t\t%d\n", server.ISO.ID)
				cmd.Printf("  Name:\t\t%s\n", server.ISO.Name)
				cmd.Printf("  Description:\t%s\n", server.ISO.Description)
				cmd.Printf("  Type:\t\t%s\n", server.ISO.Type)
			} else {
				cmd.Printf("  No ISO attached\n")
			}

			cmd.Printf("Protection:\n")
			cmd.Printf("  Delete:\t%s\n", util.YesNo(server.Protection.Delete))
			cmd.Printf("  Rebuild:\t%s\n", util.YesNo(server.Protection.Rebuild))

			cmd.Print("Labels:\n")
			if len(server.Labels) == 0 {
				cmd.Print("  No labels\n")
			} else {
				for key, value := range util.IterateInOrder(server.Labels) {
					cmd.Printf("  %s: %s\n", key, value)
				}
			}

			cmd.Print("Placement Group:\n")
			if server.PlacementGroup != nil {
				cmd.Printf("  ID:\t\t%d\n", server.PlacementGroup.ID)
				cmd.Printf("  Name:\t\t%s\n", server.PlacementGroup.Name)
				cmd.Printf("  Type:\t\t%s\n", server.PlacementGroup.Type)
			} else {
				cmd.Print("  No Placement Group set\n")
			}*/

		return nil
	},
}
