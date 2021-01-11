package floatingip

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var floatingIPListTableOutput *output.Table

func init() {
	floatingIPListTableOutput = describeFloatingIPListTableOutput(nil)
}

func newListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List Floating IPs",
		Long: util.ListLongDescription(
			"Displays a list of Floating IPs.",
			floatingIPListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runFloatingIPList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(floatingIPListTableOutput.Columns()), output.OptionJSON())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runFloatingIPList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

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

	if outOpts.IsSet("json") {
		var floatingIPSchemas []schema.FloatingIP
		for _, floatingIP := range floatingIPs {
			floatingIPSchema := schema.FloatingIP{
				ID:           floatingIP.ID,
				Name:         floatingIP.Name,
				Description:  hcloud.String(floatingIP.Description),
				IP:           floatingIP.IP.String(),
				Created:      floatingIP.Created,
				Type:         string(floatingIP.Type),
				HomeLocation: util.LocationToSchema(*floatingIP.HomeLocation),
				Blocked:      floatingIP.Blocked,
				Protection:   schema.FloatingIPProtection{Delete: floatingIP.Protection.Delete},
				Labels:       floatingIP.Labels,
			}
			for ip, dnsPtr := range floatingIP.DNSPtr {
				floatingIPSchema.DNSPtr = append(floatingIPSchema.DNSPtr, schema.FloatingIPDNSPtr{
					IP:     ip,
					DNSPtr: dnsPtr,
				})
			}
			if floatingIP.Server != nil {
				floatingIPSchema.Server = hcloud.Int(floatingIP.Server.ID)
			}
			floatingIPSchemas = append(floatingIPSchemas, floatingIPSchema)
		}
		return util.DescribeJSON(floatingIPSchemas)
	}

	cols := []string{"id", "type", "name", "description", "ip", "home", "server", "dns"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := describeFloatingIPListTableOutput(cli)
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

func describeFloatingIPListTableOutput(cli *state.State) *output.Table {
	return output.NewTable().
		AddAllowedFields(hcloud.FloatingIP{}).
		AddFieldFn("dns", output.FieldFn(func(obj interface{}) string {
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
			return util.NA(dns)
		})).
		AddFieldFn("server", output.FieldFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			var server string
			if floatingIP.Server != nil && cli != nil {
				return cli.ServerName(floatingIP.Server.ID)
			}
			return util.NA(server)
		})).
		AddFieldFn("home", output.FieldFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			return floatingIP.HomeLocation.Name
		})).
		AddFieldFn("ip", output.FieldFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			if floatingIP.Network != nil {
				return floatingIP.Network.String()
			}
			return floatingIP.IP.String()
		})).
		AddFieldFn("protection", output.FieldFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			var protection []string
			if floatingIP.Protection.Delete {
				protection = append(protection, "delete")
			}
			return strings.Join(protection, ", ")
		})).
		AddFieldFn("labels", output.FieldFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			return util.LabelsToString(floatingIP.Labels)
		})).
		AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
			floatingIP := obj.(*hcloud.FloatingIP)
			return util.Datetime(floatingIP.Created)
		}))
}
