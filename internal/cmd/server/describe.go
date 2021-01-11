package server

import (
	"encoding/json"
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newDescribeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] SERVER",
		Short:                 "Describe a server",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runServerDescribe),
	}
	output.AddFlag(cmd, output.OptionJSON(), output.OptionFormat())
	return cmd
}

func runServerDescribe(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := output.FlagsForCommand(cmd)

	idOrName := args[0]
	server, resp, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return serverDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(server, outputFlags["format"][0])
	default:
		return serverDescribeText(cli, server)
	}
}

func serverDescribeText(cli *state.State, server *hcloud.Server) error {
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
			floatingIP, _, err := cli.Client().FloatingIP.GetByID(cli.Context, f.ID)
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
			network, _, err := cli.Client().Network.GetByID(cli.Context, n.Network.ID)
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
			volume, _, err := cli.Client().Volume.GetByID(cli.Context, v.ID)
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
}

func serverDescribeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if server, ok := data["server"]; ok {
		return util.DescribeJSON(server)
	}
	if servers, ok := data["servers"].([]interface{}); ok {
		return util.DescribeJSON(servers[0])
	}
	return util.DescribeJSON(data)
}
