package cli

import (
	"strconv"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var serverListTableOutput *tableOutput

func init() {
	serverListTableOutput = describeServerListTableOutput(nil)
}

func newServerListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List servers",
		Long: listLongDescription(
			"Displays a list of servers.",
			serverListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(serverListTableOutput.Columns()), outputOptionJSON())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runServerList(cli *CLI, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	labelSelector, _ := cmd.Flags().GetString("selector")
	opts := hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
			PerPage:       50,
		},
	}
	servers, err := cli.Client().Server.AllWithOpts(cli.Context, opts)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		var serversSchema []schema.Server
		for _, server := range servers {
			serverSchema := schema.Server{
				ID:         server.ID,
				Name:       server.Name,
				Status:     string(server.Status),
				Created:    server.Created,
				Datacenter: datacenterToSchema(*server.Datacenter),
				ServerType: serverTypeToSchema(*server.ServerType),
				PublicNet: schema.ServerPublicNet{
					IPv4: schema.ServerPublicNetIPv4{
						IP:      server.PublicNet.IPv4.IP.String(),
						Blocked: server.PublicNet.IPv4.Blocked,
						DNSPtr:  server.PublicNet.IPv4.DNSPtr,
					},
					IPv6: schema.ServerPublicNetIPv6{
						IP:      server.PublicNet.IPv6.IP.String(),
						Blocked: server.PublicNet.IPv6.Blocked,
					},
				},
				RescueEnabled:   server.RescueEnabled,
				BackupWindow:    hcloud.String(server.BackupWindow),
				OutgoingTraffic: &server.OutgoingTraffic,
				IngoingTraffic:  &server.IngoingTraffic,
				IncludedTraffic: server.IncludedTraffic,
				Protection: schema.ServerProtection{
					Delete:  server.Protection.Delete,
					Rebuild: server.Protection.Rebuild,
				},
				Labels:          server.Labels,
				PrimaryDiskSize: server.PrimaryDiskSize,
			}
			if server.Image != nil {
				serverImage := imageToSchema(*server.Image)
				serverSchema.Image = &serverImage
			}
			if server.ISO != nil {
				serverISO := isoToSchema(*server.ISO)
				serverSchema.ISO = &serverISO
			}
			for ip, dnsPTR := range server.PublicNet.IPv6.DNSPtr {
				serverSchema.PublicNet.IPv6.DNSPtr = append(serverSchema.PublicNet.IPv6.DNSPtr, schema.ServerPublicNetIPv6DNSPtr{
					IP:     ip,
					DNSPtr: dnsPTR,
				})
			}
			for _, floatingIP := range server.PublicNet.FloatingIPs {
				serverSchema.PublicNet.FloatingIPs = append(serverSchema.PublicNet.FloatingIPs, floatingIP.ID)
			}
			for _, volume := range server.Volumes {
				serverSchema.Volumes = append(serverSchema.Volumes, volume.ID)
			}
			for _, privateNet := range server.PrivateNet {
				privateNetSchema := schema.ServerPrivateNet{
					Network:    privateNet.Network.ID,
					IP:         privateNet.IP.String(),
					MACAddress: privateNet.MACAddress,
				}
				for _, aliasIP := range privateNet.Aliases {
					privateNetSchema.AliasIPs = append(privateNetSchema.AliasIPs, aliasIP.String())
				}
				serverSchema.PrivateNet = append(serverSchema.PrivateNet, privateNetSchema)
			}
			serversSchema = append(serversSchema, serverSchema)
		}
		return describeJSON(serversSchema)
	}

	cols := []string{"id", "name", "status", "ipv4", "ipv6", "datacenter"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := describeServerListTableOutput(cli)
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, server := range servers {
		tw.Write(cols, server)
	}
	tw.Flush()
	return nil
}

func describeServerListTableOutput(cli *CLI) *tableOutput {
	return newTableOutput().
		AddAllowedFields(hcloud.Server{}).
		AddFieldOutputFn("ipv4", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return server.PublicNet.IPv4.IP.String()
		})).
		AddFieldOutputFn("ipv6", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return server.PublicNet.IPv6.Network.String()
		})).
		AddFieldOutputFn("datacenter", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return server.Datacenter.Name
		})).
		AddFieldOutputFn("location", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return server.Datacenter.Location.Name
		})).
		AddFieldOutputFn("labels", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return labelsToString(server.Labels)
		})).
		AddFieldOutputFn("type", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return server.ServerType.Name
		})).
		AddFieldOutputFn("volumes", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			var volumes []string
			for _, volume := range server.Volumes {
				volumeID := strconv.Itoa(volume.ID)
				volumes = append(volumes, volumeID)
			}
			return strings.Join(volumes, ", ")
		})).
		AddFieldOutputFn("private_net", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			var networks []string
			if cli != nil {
				for _, network := range server.PrivateNet {
					networks = append(networks, cli.NetworkName(network.Network.ID))
				}
			}
			return na(strings.Join(networks, ", "))
		})).
		AddFieldOutputFn("protection", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			var protection []string
			if server.Protection.Delete {
				protection = append(protection, "delete")
			}
			if server.Protection.Rebuild {
				protection = append(protection, "rebuild")
			}
			return strings.Join(protection, ", ")
		})).
		AddFieldOutputFn("created", fieldOutputFn(func(obj interface{}) string {
			server := obj.(*hcloud.Server)
			return datetime(server.Created)
		}))
}
