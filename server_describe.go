package cli

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func newServerDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] SERVER",
		Short:                 "Describe a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerDescribe),
	}
	return cmd
}

func runServerDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	fmt.Printf("ID:\t\t%d\n", server.ID)
	fmt.Printf("Name:\t\t%s\n", server.Name)
	fmt.Printf("Status:\t\t%s\n", server.Status)
	fmt.Printf("Created:\t%s (%s)\n", datetime(server.Created), humanize.Time(server.Created))

	fmt.Printf("Server Type:\t%s (ID: %d)\n", server.ServerType.Name, server.ServerType.ID)
	fmt.Printf("  ID:\t\t%d\n", server.ServerType.ID)
	fmt.Printf("  Name:\t\t%s\n", server.ServerType.Name)
	fmt.Printf("  Description:\t%s\n", server.ServerType.Description)
	fmt.Printf("  Cores:\t%d\n", server.ServerType.Cores)
	fmt.Printf("  Memory:\t%v GB\n", server.ServerType.Memory)
	fmt.Printf("  Disk:\t\t%d GB\n", server.ServerType.Disk)
	fmt.Printf("  Storage Type:\t%s\n", server.ServerType.StorageType)

	fmt.Printf("Public Net:\n")
	fmt.Printf("  IPv4:\n")
	fmt.Printf("    IP:\t\t%s\n", server.PublicNet.IPv4.IP)
	fmt.Printf("    Blocked:\t%s\n", yesno(server.PublicNet.IPv4.Blocked))
	fmt.Printf("    DNS:\t%s\n", server.PublicNet.IPv4.DNSPtr)
	fmt.Printf("  IPv6:\n")
	fmt.Printf("    IP:\t\t%s\n", server.PublicNet.IPv6.Network.String())
	fmt.Printf("    Blocked:\t%s\n", yesno(server.PublicNet.IPv6.Blocked))
	fmt.Printf("  Floating IPs:\n")
	if len(server.PublicNet.FloatingIPs) > 0 {
		for _, f := range server.PublicNet.FloatingIPs {
			floatingIP, _, err := cli.client.FloatingIP.GetByID(cli.Context, f.ID)
			if err != nil {
				return fmt.Errorf("error fetching Floating IP: %v", err)
			}
			fmt.Printf("  - ID:\t\t\t%d\n", floatingIP.ID)
			fmt.Printf("    Description:\t%s\n", na(floatingIP.Description))
			fmt.Printf("    IP:\t\t\t%s\n", floatingIP.IP)
		}
	} else {
		fmt.Printf("    No Floating IPs\n")
	}
	fmt.Printf("Image:\n")
	if server.Image != nil {
		image := server.Image
		fmt.Printf("  ID:\t\t%d\n", image.ID)
		fmt.Printf("  Type:\t\t%s\n", image.Type)
		fmt.Printf("  Status:\t%s\n", image.Status)
		fmt.Printf("  Name:\t\t%s\n", na(image.Name))
		fmt.Printf("  Description:\t%s\n", image.Description)
		if image.ImageSize != 0 {
			fmt.Printf("  Image size:\t%.1f GB\n", image.ImageSize)
		} else {
			fmt.Printf("  Image size:\t%s\n", na(""))
		}
		fmt.Printf("  Disk size:\t%.0f GB\n", image.DiskSize)
		fmt.Printf("  Created:\t%s (%s)\n", datetime(image.Created), humanize.Time(image.Created))
		fmt.Printf("  OS flavor:\t%s\n", image.OSFlavor)
		fmt.Printf("  OS version:\t%s\n", na(image.OSVersion))
		fmt.Printf("  Rapid deploy:\t%s\n", yesno(image.RapidDeploy))
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
	fmt.Printf("  Outgoing:\t%v\n", humanize.Bytes(server.OutgoingTraffic))
	fmt.Printf("  Ingoing:\t%v\n", humanize.Bytes(server.IngoingTraffic))
	fmt.Printf("  Included:\t%v\n", humanize.Bytes(server.IncludedTraffic))

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

	return nil
}
