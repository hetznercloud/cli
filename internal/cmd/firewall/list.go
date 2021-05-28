package firewall

import (
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var listTableOutput *output.Table

func init() {
	listTableOutput = output.NewTable().
		AddAllowedFields(hcloud.Firewall{})
}

func newListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List Firewalls",
		Long: util.ListLongDescription(
			"Displays a list of Firewalls.",
			listTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(listTableOutput.Columns()), output.OptionJSON())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)
	labelSelector, _ := cmd.Flags().GetString("selector")
	opts := hcloud.FirewallListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
			PerPage:       50,
		},
	}
	firewalls, err := cli.Client().Firewall.AllWithOpts(cli.Context, opts)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		var firewallSchemas []schema.Firewall
		for _, firewall := range firewalls {
			firewallSchema := schema.Firewall{
				ID:      firewall.ID,
				Name:    firewall.Name,
				Labels:  firewall.Labels,
				Created: firewall.Created,
			}
			for _, rule := range firewall.Rules {
				var sourceNets []string
				for _, sourceIP := range rule.SourceIPs {
					sourceNets = append(sourceNets, sourceIP.Network())
				}
				var destinationNets []string
				for _, destinationIP := range rule.DestinationIPs {
					destinationNets = append(destinationNets, destinationIP.Network())
				}
				firewallSchema.Rules = append(firewallSchema.Rules, schema.FirewallRule{
					Direction:      string(rule.Direction),
					SourceIPs:      sourceNets,
					DestinationIPs: destinationNets,
					Protocol:       string(rule.Protocol),
					Port:           rule.Port,
				})
			}
			for _, AppliedTo := range firewall.AppliedTo {
				s := schema.FirewallResource{
					Type: string(AppliedTo.Type),
				}
				switch AppliedTo.Type {
				case hcloud.FirewallResourceTypeServer:
					s.Server = &schema.FirewallResourceServer{ID: AppliedTo.Server.ID}
				}
				firewallSchema.AppliedTo = append(firewallSchema.AppliedTo, s)
			}

			firewallSchemas = append(firewallSchemas, firewallSchema)
		}
		return util.DescribeJSON(firewallSchemas)
	}

	cols := []string{"id", "name"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := listTableOutput
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, firewall := range firewalls {
		tw.Write(cols, firewall)
	}
	tw.Flush()

	return nil
}
