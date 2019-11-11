package cli

import (
	"encoding/json"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newFloatingIPDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] FLOATINGIP",
		Short:                 "Describe a Floating IP",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runFloatingIPDescribe),
	}
	addOutputFlag(cmd, outputOptionJSON(), outputOptionFormat())
	return cmd
}

func runFloatingIPDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	outputFlags := outputFlagsForCommand(cmd)

	idOrName := args[0]
	floatingIP, resp, err := cli.Client().FloatingIP.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if floatingIP == nil {
		return fmt.Errorf("Floating IP not found: %v", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return floatingIPDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return describeFormat(floatingIP, outputFlags["format"][0])
	default:
		return floatingIPDescribeText(cli, floatingIP)
	}
}

func floatingIPDescribeText(cli *CLI, floatingIP *hcloud.FloatingIP) error {
	fmt.Printf("ID:\t\t%d\n", floatingIP.ID)
	fmt.Printf("Type:\t\t%s\n", floatingIP.Type)
	fmt.Printf("Name:\t\t%s\n", floatingIP.Name)
	fmt.Printf("Description:\t%s\n", na(floatingIP.Description))
	fmt.Printf("Created:\t%s (%s)\n", datetime(floatingIP.Created), humanize.Time(floatingIP.Created))
	if floatingIP.Network != nil {
		fmt.Printf("IP:\t\t%s\n", floatingIP.Network.String())
	} else {
		fmt.Printf("IP:\t\t%s\n", floatingIP.IP.String())
	}
	fmt.Printf("Blocked:\t%s\n", yesno(floatingIP.Blocked))
	fmt.Printf("Home Location:\t%s\n", floatingIP.HomeLocation.Name)
	if floatingIP.Server != nil {
		server, _, err := cli.Client().Server.GetByID(cli.Context, floatingIP.Server.ID)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %d", floatingIP.Server.ID)
		}
		fmt.Printf("Server:\n")
		fmt.Printf("  ID:\t%d\n", server.ID)
		fmt.Printf("  Name:\t%s\n", server.Name)
	} else {
		fmt.Print("Server:\n  Not assigned\n")
	}
	fmt.Print("DNS:\n")
	if len(floatingIP.DNSPtr) == 0 {
		fmt.Print("  No reverse DNS entries\n")
	} else {
		for ip, dns := range floatingIP.DNSPtr {
			fmt.Printf("  %s: %s\n", ip, dns)
		}
	}

	fmt.Printf("Protection:\n")
	fmt.Printf("  Delete:\t%s\n", yesno(floatingIP.Protection.Delete))

	fmt.Print("Labels:\n")
	if len(floatingIP.Labels) == 0 {
		fmt.Print("  No labels\n")
	} else {
		for key, value := range floatingIP.Labels {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
	return nil
}

func floatingIPDescribeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if floatingIP, ok := data["floating_ip"]; ok {
		return describeJSON(floatingIP)
	}
	if floatingIPs, ok := data["floating_ips"].([]interface{}); ok {
		return describeJSON(floatingIPs[0])
	}
	return describeJSON(data)
}
