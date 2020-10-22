package cli

import (
	"encoding/json"
	"fmt"
	"github.com/guptarohit/asciigraph"
	"sort"
	"strconv"
	"time"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerMetricsCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "metrics [FLAGS] SERVER",
		Short:                 "[ALPHA] Metrics from a Server",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerMetrics),
	}

	cmd.Flags().String("type", "", "Fancy Graphs")
	cmd.MarkFlagRequired("type")
	cmd.Flags().String("start", "", "ISO 8601 timestamp")
	cmd.Flags().String("end", "", "ISO 8601 timestamp")

	addOutputFlag(cmd, outputOptionJSON())
	return cmd
}

func runServerMetrics(cli *CLI, cmd *cobra.Command, args []string) error {
	outputFlags := outputFlagsForCommand(cmd)

	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
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

	m, resp, err := cli.Client().Server.GetMetrics(cli.Context, server, hcloud.ServerGetMetricsOpts{
		Types: []hcloud.ServerMetricType{hcloud.ServerMetricType(metricType)},
		Start: startTime,
		End:   endTime,
	})
	if err != nil {
		return err
	}
	switch {
	case outputFlags.IsSet("json"):
		return serverMetricsJSON(resp)
	default:
		var keys []string
		for k := range m.TimeSeries {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if len(m.TimeSeries[k]) == 0 {
				fmt.Printf("Currently there are now metrics available. Please try it again later.")
				return nil
			}
			fmt.Printf("Server: %s \t Metric: %s \t Start: %s \t End: %s\n", server.Name, k, m.Start.String(), m.End.String())
			var data []float64
			for _, m := range m.TimeSeries[k] {
				d, _ := strconv.ParseFloat(m.Value, 64)
				data = append(data, d)
			}
			graph := asciigraph.Plot(data, asciigraph.Height(20), asciigraph.Width(100))
			fmt.Println(graph)
			fmt.Printf("\n\n")
		}
	}
	return nil
}

func serverMetricsJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	return describeJSON(data)
}
