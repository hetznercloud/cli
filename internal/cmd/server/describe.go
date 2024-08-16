package server

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "server",
	ShortDescription:     "Describe a server",
	JSONKeyGetByID:       "server",
	JSONKeyGetByName:     "servers",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Server().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, interface{}, error) {
		srv, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return srv, hcloud.SchemaFromServer(srv), nil
	},
	PrintText: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		server := resource.(*hcloud.Server)

		cmd.Printf("ID:\t\t%d\n", server.ID)
		cmd.Printf("Name:\t\t%s\n", server.Name)
		cmd.Printf("Status:\t\t%s\n", server.Status)
		cmd.Printf("Created:\t%s (%s)\n", util.Datetime(server.Created), humanize.Time(server.Created))

		cmd.Printf("Server Type:\t%s (ID: %d)\n", server.ServerType.Name, server.ServerType.ID)
		cmd.Printf("  ID:\t\t%d\n", server.ServerType.ID)
		cmd.Printf("  Name:\t\t%s\n", server.ServerType.Name)
		cmd.Printf("  Description:\t%s\n", server.ServerType.Description)
		cmd.Printf("  Cores:\t%d\n", server.ServerType.Cores)
		cmd.Printf("  CPU Type:\t%s\n", server.ServerType.CPUType)
		cmd.Printf("  Memory:\t%v GB\n", server.ServerType.Memory)
		cmd.Printf("  Disk:\t\t%d GB\n", server.PrimaryDiskSize)
		cmd.Printf("  Storage Type:\t%s\n", server.ServerType.StorageType)
		cmd.Print(util.PrefixLines(util.DescribeDeprecation(server.ServerType), "  "))

		cmd.Printf("Public Net:\n")
		cmd.Printf("  IPv4:\n")
		if server.PublicNet.IPv4.IsUnspecified() {
			cmd.Printf("    No Primary IPv4\n")
		} else {
			cmd.Printf("    ID:\t\t%d\n", server.PublicNet.IPv4.ID)
			cmd.Printf("    IP:\t\t%s\n", server.PublicNet.IPv4.IP)
			cmd.Printf("    Blocked:\t%s\n", util.YesNo(server.PublicNet.IPv4.Blocked))
			cmd.Printf("    DNS:\t%s\n", server.PublicNet.IPv4.DNSPtr)
		}

		cmd.Printf("  IPv6:\n")
		if server.PublicNet.IPv6.IsUnspecified() {
			cmd.Printf("    No Primary IPv6\n")
		} else {
			cmd.Printf("    ID:\t\t%d\n", server.PublicNet.IPv6.ID)
			cmd.Printf("    IP:\t\t%s\n", server.PublicNet.IPv6.Network.String())
			cmd.Printf("    Blocked:\t%s\n", util.YesNo(server.PublicNet.IPv6.Blocked))
		}
		cmd.Printf("  Floating IPs:\n")
		if len(server.PublicNet.FloatingIPs) > 0 {
			for _, f := range server.PublicNet.FloatingIPs {
				floatingIP, _, err := s.Client().FloatingIP().GetByID(s, f.ID)
				if err != nil {
					return fmt.Errorf("error fetching Floating IP: %v", err)
				}
				cmd.Printf("  - ID:\t\t\t%d\n", floatingIP.ID)
				cmd.Printf("    Description:\t%s\n", util.NA(floatingIP.Description))
				cmd.Printf("    IP:\t\t\t%s\n", floatingIP.IP)
			}
		} else {
			cmd.Printf("    No Floating IPs\n")
		}

		cmd.Printf("Private Net:\n")
		if len(server.PrivateNet) > 0 {
			for _, n := range server.PrivateNet {
				network, _, err := s.Client().Network().GetByID(s, n.Network.ID)
				if err != nil {
					return fmt.Errorf("error fetching network: %v", err)
				}
				cmd.Printf("  - ID:\t\t\t%d\n", network.ID)
				cmd.Printf("    Name:\t\t%s\n", network.Name)
				cmd.Printf("    IP:\t\t\t%s\n", n.IP.String())
				cmd.Printf("    MAC Address:\t%s\n", n.MACAddress)
				if len(n.Aliases) > 0 {
					cmd.Printf("    Alias IPs:\n")
					for _, a := range n.Aliases {
						cmd.Printf("     -\t\t\t%s\n", a)
					}
				} else {
					cmd.Printf("    Alias IPs:\t\t%s\n", util.NA(""))
				}
			}
		} else {
			cmd.Printf("    No Private Networks\n")
		}

		cmd.Printf("Volumes:\n")
		if len(server.Volumes) > 0 {
			for _, v := range server.Volumes {
				volume, _, err := s.Client().Volume().GetByID(s, v.ID)
				if err != nil {
					return fmt.Errorf("error fetching Volume: %v", err)
				}
				cmd.Printf("  - ID:\t\t%d\n", volume.ID)
				cmd.Printf("    Name:\t%s\n", volume.Name)
				cmd.Printf("    Size:\t%s\n", humanize.Bytes(uint64(volume.Size)*humanize.GByte))
			}
		} else {
			cmd.Printf("  No Volumes\n")
		}
		cmd.Printf("Image:\n")
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
			for key, value := range server.Labels {
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
		}

		return nil
	},
}
