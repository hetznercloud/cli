package cmds

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/guptarohit/asciigraph"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerMetricsCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "metrics [FLAGS] LOADBALANCER",
		Short:                 "[ALPHA] Metrics from a Load Balancer",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.LoadBalancerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runLoadBalancerMetrics),
	}

	cmd.Flags().String("type", "", "Type of metrics you want to show")
	cmd.MarkFlagRequired("type")
	cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("open_connections", "connections_per_second", "requests_per_second", "bandwidth"))

	cmd.Flags().String("start", "", "ISO 8601 timestamp")
	cmd.Flags().String("end", "", "ISO 8601 timestamp")

	addOutputFlag(cmd, outputOptionJSON())
	return cmd
}

func runLoadBalancerMetrics(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := outputFlagsForCommand(cmd)

	idOrName := args[0]
	LoadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if LoadBalancer == nil {
		return fmt.Errorf("LoadBalancer not found: %s", idOrName)
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

	m, resp, err := cli.Client().LoadBalancer.GetMetrics(cli.Context, LoadBalancer, hcloud.LoadBalancerGetMetricsOpts{
		Types: []hcloud.LoadBalancerMetricType{hcloud.LoadBalancerMetricType(metricType)},
		Start: startTime,
		End:   endTime,
	})

	if err != nil {
		return err
	}
	switch {
	case outputFlags.IsSet("json"):
		return loadBalancerMetricsJSON(resp)
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
			fmt.Printf("Load Balancer: %s \t Metric: %s \t Start: %s \t End: %s\n", LoadBalancer.Name, k, m.Start.String(), m.End.String())
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

func loadBalancerMetricsJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	return describeJSON(data)
}
