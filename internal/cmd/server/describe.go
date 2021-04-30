package server

import (
	"context"
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var describeCmd = base.DescribeCmd{
	ResourceNameSingular: "server",
	ShortDescription:     "Describe a server",
	JSONKeyGetByID:       "server",
	JSONKeyGetByName:     "servers",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Server().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Server().Get(ctx, idOrName)
	},
	PrintText: func(ctx context.Context, client hcapi2.Client, resource interface{}) error {
		server := resource.(*hcloud.Server)

		fmt.Printf("ID:\t\t%d\n", server.ID)
		fmt.Printf("Name:\t\t%s\n", server.Name)
		fmt.Printf("Status:\t\t%s\n", server.Status)
		fmt.Printf("Created:\t%s (%s)\n", util.Datetime(server.Created), humanize.Time(server.Created))

		fmt.Printf("Server Type:\t%s (ID: %d)\n", server.ServerType.Name, server.ServerType.ID)
		fmt.Printf("  ID:\t\t%d\n", server.ServerType.ID)
		fmt.Printf("  Name:\t\t%s\n", server.ServerType.Name)
		fmt.Printf("  Description:\t%s\n", server.ServerType.Description)
		fmt.Printf("  Cores:\t%d\n", server.ServerType.Cores)
		fmt.Printf("  CPU Type:\t%s\n", server.ServerType.CPUType)
		fmt.Printf("  Memory:\t%v GB\n", server.ServerType.Memory)
		fmt.Printf("  Disk:\t\t%d GB\n", server.PrimaryDiskSize)
		fmt.Printf("  Storage Type:\t%s\n", server.ServerType.StorageType)

		fmt.Printf("Public Net:\n")
		fmt.Printf("  IPv4:\n")
		fmt.Printf("    IP:\t\t%s\n", server.PublicNet.IPv4.IP)
		fmt.Printf("    Blocked:\t%s\n", util.YesNo(server.PublicNet.IPv4.Blocked))
		fmt.Printf("    DNS:\t%s\n", server.PublicNet.IPv4.DNSPtr)
		fmt.Printf("  IPv6:\n")
		fmt.Printf("    IP:\t\t%s\n", server.PublicNet.IPv6.Network.String())
		fmt.Printf("    Blocked:\t%s\n", util.YesNo(server.PublicNet.IPv6.Blocked))
		fmt.Printf("  Floating IPs:\n")
		if len(server.PublicNet.FloatingIPs) > 0 {
			for _, f := range server.PublicNet.FloatingIPs {
				floatingIP, _, err := client.FloatingIP().GetByID(ctx, f.ID)
				if err != nil {
					return fmt.Errorf("error fetching Floating IP: %v", err)
				}
				fmt.Printf("  - ID:\t\t\t%d\n", floatingIP.ID)
				fmt.Printf("    Description:\t%s\n", util.NA(floatingIP.Description))
				fmt.Printf("    IP:\t\t\t%s\n", floatingIP.IP)
			}
		} else {
			fmt.Printf("    No Floating IPs\n")
		}

		fmt.Printf("Private Net:\n")
		if len(server.PrivateNet) > 0 {
			for _, n := range server.PrivateNet {
				network, _, err := client.Network().GetByID(ctx, n.Network.ID)
				if err != nil {
					return fmt.Errorf("error fetching network: %v", err)
				}
				fmt.Printf("  - ID:\t\t\t%d\n", network.ID)
				fmt.Printf("    Name:\t\t%s\n", network.Name)
				fmt.Printf("    IP:\t\t\t%s\n", n.IP.String())
				fmt.Printf("    MAC Address:\t%s\n", n.MACAddress)
				if len(n.Aliases) > 0 {
					fmt.Printf("    Alias IPs:\n")
					for _, a := range n.Aliases {
						fmt.Printf("     -\t\t\t%s\n", a)
					}
				} else {
					fmt.Printf("    Alias IPs:\t\t%s\n", util.NA(""))
				}
			}
		} else {
			fmt.Printf("    No Private Networks\n")
		}

		fmt.Printf("Volumes:\n")
		if len(server.Volumes) > 0 {
			for _, v := range server.Volumes {
				volume, _, err := client.Volume().GetByID(ctx, v.ID)
				if err != nil {
					return fmt.Errorf("error fetching Volume: %v", err)
				}
				fmt.Printf("  - ID:\t\t%d\n", volume.ID)
				fmt.Printf("    Name:\t%s\n", volume.Name)
				fmt.Printf("    Size:\t%s\n", humanize.Bytes(uint64(volume.Size*humanize.GByte)))
			}
		} else {
			fmt.Printf("  No Volumes\n")
		}
		fmt.Printf("Image:\n")
		if server.Image != nil {
			image := server.Image
			fmt.Printf("  ID:\t\t%d\n", image.ID)
			fmt.Printf("  Type:\t\t%s\n", image.Type)
			fmt.Printf("  Status:\t%s\n", image.Status)
			fmt.Printf("  Name:\t\t%s\n", util.NA(image.Name))
			fmt.Printf("  Description:\t%s\n", image.Description)
			if image.ImageSize != 0 {
				fmt.Printf("  Image size:\t%.2f GB\n", image.ImageSize)
			} else {
				fmt.Printf("  Image size:\t%s\n", util.NA(""))
			}
			fmt.Printf("  Disk size:\t%.0f GB\n", image.DiskSize)
			fmt.Printf("  Created:\t%s (%s)\n", util.Datetime(image.Created), humanize.Time(image.Created))
			fmt.Printf("  OS flavor:\t%s\n", image.OSFlavor)
			fmt.Printf("  OS version:\t%s\n", util.NA(image.OSVersion))
			fmt.Printf("  Rapid deploy:\t%s\n", util.YesNo(image.RapidDeploy))
		} else {
			fmt.Printf("  No Image\n")
		}

		fmt.Printf("Datacenter:\n")
		fmt.Printf("  ID:\t\t%d\n", server.Datacenter.ID)
		fmt.Printf("  Name:\t\t%s\n", server.Datacenter.Name)
		fmt.Printf("  Description:\t%s\n", server.Datacenter.Description)
		fmt.Printf("  Location:\n")
		fmt.Printf("    Name:\t\t%s\n", server.Datacenter.Location.Name)
		fmt.Printf("    Description:\t%s\n", server.Datacenter.Location.Description)
		fmt.Printf("    Country:\t\t%s\n", server.Datacenter.Location.Country)
		fmt.Printf("    City:\t\t%s\n", server.Datacenter.Location.City)
		fmt.Printf("    Latitude:\t\t%f\n", server.Datacenter.Location.Latitude)
		fmt.Printf("    Longitude:\t\t%f\n", server.Datacenter.Location.Longitude)

		fmt.Printf("Traffic:\n")
		fmt.Printf("  Outgoing:\t%v\n", humanize.IBytes(server.OutgoingTraffic))
		fmt.Printf("  Ingoing:\t%v\n", humanize.IBytes(server.IngoingTraffic))
		fmt.Printf("  Included:\t%v\n", humanize.IBytes(server.IncludedTraffic))

		if server.BackupWindow != "" {
			fmt.Printf("Backup Window:\t%s\n", server.BackupWindow)
		} else {
			fmt.Printf("Backup Window:\tBackups disabled\n")
		}

		if server.RescueEnabled {
			fmt.Printf("Rescue System:\tenabled\n")
		} else {
			fmt.Printf("Rescue System:\tdisabled\n")
		}

		fmt.Printf("ISO:\n")
		if server.ISO != nil {
			fmt.Printf("  ID:\t\t%d\n", server.ISO.ID)
			fmt.Printf("  Name:\t\t%s\n", server.ISO.Name)
			fmt.Printf("  Description:\t%s\n", server.ISO.Description)
			fmt.Printf("  Type:\t\t%s\n", server.ISO.Type)
		} else {
			fmt.Printf("  No ISO attached\n")
		}

		fmt.Printf("Protection:\n")
		fmt.Printf("  Delete:\t%s\n", util.YesNo(server.Protection.Delete))
		fmt.Printf("  Rebuild:\t%s\n", util.YesNo(server.Protection.Rebuild))

		fmt.Print("Labels:\n")
		if len(server.Labels) == 0 {
			fmt.Print("  No labels\n")
		} else {
			for key, value := range server.Labels {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}

		return nil
	},
}
