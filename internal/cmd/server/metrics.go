package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var MetricsCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "metrics [FLAGS] SERVER",
			Short:                 "[ALPHA] Metrics from a Server",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("type", "", "Type of metrics you want to show")
		cmd.MarkFlagRequired("type")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("cpu", "disk", "network"))

		cmd.Flags().String("start", "", "ISO 8601 timestamp")
		cmd.Flags().String("end", "", "ISO 8601 timestamp")

		output.AddFlag(cmd, output.OptionJSON())
		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		outputFlags := output.FlagsForCommand(cmd)

		idOrName := args[0]
		server, _, err := client.Server().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		metricType, _ := cmd.Flags().GetString("type")

		start, _ := cmd.Flags().GetString("start")
		startTime := time.Now().Add(-30 * time.Minute)
		if start != "" {
			startTime, err = time.Parse(time.RFC3339, start)
			if err != nil {
				return fmt.Errorf("start date has an invalid format. It should be ISO 8601, like: %s", time.Now().Format(time.RFC3339))
			}
		}

		end, _ := cmd.Flags().GetString("end")
		endTime := time.Now()
		if end != "" {
			endTime, err = time.Parse(time.RFC3339, end)
			if err != nil {
				return fmt.Errorf("end date has an invalid format. It should be ISO 8601, like: %s", time.Now().Format(time.RFC3339))
			}
		}

		m, resp, err := client.Server().GetMetrics(ctx, server, hcloud.ServerGetMetricsOpts{
			Types: []hcloud.ServerMetricType{hcloud.ServerMetricType(metricType)},
			Start: startTime,
			End:   endTime,
		})
		if err != nil {
			return err
		}
		switch {
		case outputFlags.IsSet("json"):
			var data map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				return err
			}
			return util.DescribeJSON(data)
		default:
			var keys []string
			for k := range m.TimeSeries {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				if len(m.TimeSeries[k]) == 0 {
					cmd.Printf("Currently there are now metrics available. Please try it again later.")
					return nil
				}
				cmd.Printf("Server: %s \t Metric: %s \t Start: %s \t End: %s\n", server.Name, k, m.Start.String(), m.End.String())
				var data []float64
				for _, m := range m.TimeSeries[k] {
					d, _ := strconv.ParseFloat(m.Value, 64)
					data = append(data, d)
				}
				graph := asciigraph.Plot(data, asciigraph.Height(20), asciigraph.Width(100))
				cmd.Println(graph)
				cmd.Printf("\n\n")
			}
		}
		return nil
	},
}
