package all

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/certificate"
	"github.com/hetznercloud/cli/internal/cmd/firewall"
	"github.com/hetznercloud/cli/internal/cmd/floatingip"
	"github.com/hetznercloud/cli/internal/cmd/image"
	"github.com/hetznercloud/cli/internal/cmd/iso"
	"github.com/hetznercloud/cli/internal/cmd/loadbalancer"
	"github.com/hetznercloud/cli/internal/cmd/network"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/placementgroup"
	"github.com/hetznercloud/cli/internal/cmd/primaryip"
	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/cmd/sshkey"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/cmd/volume"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var allCmds = []base.ListCmd{
	server.ListCmd,
	image.ListCmd,
	placementgroup.ListCmd,
	primaryip.ListCmd,
	iso.ListCmd,
	volume.ListCmd,
	loadbalancer.ListCmd,
	floatingip.ListCmd,
	network.ListCmd,
	firewall.ListCmd,
	certificate.ListCmd,
	sshkey.ListCmd,
}

var ListCmd = base.Cmd{
	BaseCobraCommand: func(hcapi2.Client) *cobra.Command {

		var resources []string
		for _, cmd := range allCmds {
			resources = append(resources, " - "+cmd.ResourceNamePlural)
		}

		cmd := &cobra.Command{
			Use:   "list [options]",
			Short: "List all resources in the project",
			Long: `List all resources in the project. This does not include static/public resources like locations, public ISOs, etc.

Listed resources are:
` + strings.Join(resources, "\n"),
		}

		cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")

		cmd.Flags().Bool("paid", false, "Only list resources that cost money")

		output.AddFlag(cmd, output.OptionJSON(), output.OptionYAML())

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, _ []string) error {

		paid, _ := cmd.Flags().GetBool("paid")
		labelSelector, _ := cmd.Flags().GetString("selector")

		outOpts := output.FlagsForCommand(cmd)

		var cmds []base.ListCmd
		if paid {
			cmds = []base.ListCmd{
				server.ListCmd,
				image.ListCmd,
				primaryip.ListCmd,
				volume.ListCmd,
				loadbalancer.ListCmd,
				floatingip.ListCmd,
			}
		} else {
			cmds = []base.ListCmd{
				server.ListCmd,
				image.ListCmd,
				placementgroup.ListCmd,
				primaryip.ListCmd,
				iso.ListCmd,
				volume.ListCmd,
				loadbalancer.ListCmd,
				floatingip.ListCmd,
				network.ListCmd,
				firewall.ListCmd,
				certificate.ListCmd,
				sshkey.ListCmd,
			}
		}

		type response struct {
			result []any
			err    error
		}
		responseChs := make([]chan response, len(cmds))

		// Start all requests in parallel in order to minimize response time
		for i, lc := range cmds {
			i, lc := i, lc
			ch := make(chan response)
			responseChs[i] = ch

			go func() {
				defer close(ch)

				listOpts := hcloud.ListOpts{
					LabelSelector: labelSelector,
				}

				flagSet := pflag.NewFlagSet(lc.JSONKeyGetByName, pflag.ExitOnError)

				switch lc.JSONKeyGetByName {
				case image.ListCmd.JSONKeyGetByName:
					flagSet.StringSlice("type", []string{"backup", "snapshot"}, "")
				case iso.ListCmd.JSONKeyGetByName:
					flagSet.StringSlice("type", []string{"private"}, "")
				}

				// FlagSet has to be parsed to be populated.
				// We pass an empty slice because we defined the flags earlier.
				_ = flagSet.Parse([]string{})

				result, err := lc.Fetch(s, flagSet, listOpts, []string{})
				ch <- response{result, err}
			}()
		}

		// Wait for all requests to finish and collect results
		resources := make([][]any, len(cmds))
		for i, responseCh := range responseChs {
			response := <-responseCh
			if err := response.err; err != nil {
				return err
			}
			resources[i] = response.result
		}

		if outOpts.IsSet("json") || outOpts.IsSet("yaml") {
			schema := make(map[string]any)
			for i, lc := range cmds {
				schema[lc.JSONKeyGetByName] = lc.Schema(resources[i])
			}
			if outOpts.IsSet("json") {
				return util.DescribeJSON(schema)
			}
			return util.DescribeYAML(schema)
		}

		for i, lc := range cmds {
			cols := lc.DefaultColumns
			table := lc.OutputTable(s.Client())
			table.WriteHeader(cols)

			if len(resources[i]) == 0 {
				continue
			}

			cmd.Print(strings.ToUpper(lc.ResourceNamePlural) + "\n---\n")
			for _, resource := range resources[i] {
				table.Write(cols, resource)
			}
			if err := table.Flush(); err != nil {
				return err
			}
			cmd.Println()
		}

		return nil
	},
}
