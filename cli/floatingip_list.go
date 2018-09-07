package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var floatingIPListTableOutput *tableOutput

func init() {
	floatingIPListTableOutput = newTableOutput().
		AddAllowedFields(hcloud.FloatingIP{}).
		AddFieldOutputFn("dns", fieldOutputFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			var dns string
			if len(floatingIP.DNSPtr) == 1 {
				for _, v := range floatingIP.DNSPtr {
					dns = v
				}
			}
			if len(floatingIP.DNSPtr) > 1 {
				dns = fmt.Sprintf("%d entries", len(floatingIP.DNSPtr))
			}
			return na(dns)
		})).
		AddFieldOutputFn("server", fieldOutputFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			var server string
			if floatingIP.Server != nil {
				server = strconv.Itoa(floatingIP.Server.ID)
			}
			return na(server)
		})).
		AddFieldOutputFn("home", fieldOutputFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			return floatingIP.HomeLocation.Name
		})).
		AddFieldOutputFn("ip", fieldOutputFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			if floatingIP.Network != nil {
				return floatingIP.Network.String()
			}
			return floatingIP.IP.String()
		})).
		AddFieldOutputFn("protection", fieldOutputFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			var protection []string
			if floatingIP.Protection.Delete {
				protection = append(protection, "delete")
			}
			return strings.Join(protection, ", ")
		}))
}

func newFloatingIPListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List Floating IPs",
		Long: listLongDescription(
			"Displays a list of Floating IPs.",
			floatingIPListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runFloatingIPList),
	}
	addListOutputFlag(cmd, floatingIPListTableOutput.Columns())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runFloatingIPList(cli *CLI, cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetStringArray("output")
	outOpts, err := parseOutputOpts(out)
	if err != nil {
		return err
	}

	labelSelector, _ := cmd.Flags().GetString("selector")
	opts := hcloud.FloatingIPListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
			PerPage:       50,
		},
	}
	floatingIPs, err := cli.Client().FloatingIP.AllWithOpts(cli.Context, opts)
	if err != nil {
		return err
	}

	cols := []string{"id", "type", "description", "ip", "home", "server", "dns"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := floatingIPListTableOutput
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, floatingIP := range floatingIPs {
		tw.Write(cols, floatingIP)
	}
	tw.Flush()
	return nil
}
