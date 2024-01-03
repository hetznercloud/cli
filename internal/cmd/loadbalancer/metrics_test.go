package loadbalancer

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestMetrics(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := MetricsCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	start := time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		GetMetrics(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerGetMetricsOpts{
			Start: start,
			End:   end,
			Types: []hcloud.LoadBalancerMetricType{hcloud.LoadBalancerMetricOpenConnections},
		}).
		Return(&hcloud.LoadBalancerMetrics{
			TimeSeries: map[string][]hcloud.LoadBalancerMetricsValue{
				"open_connections": {
					{
						Timestamp: float64(start.Add(0*time.Minute).UnixMilli() / 1000.0),
						Value:     "2",
					},
					{
						Timestamp: float64(start.Add(10*time.Minute).UnixMilli() / 1000.0),
						Value:     "4",
					},
					{
						Timestamp: float64(start.Add(20*time.Minute).UnixMilli() / 1000.0),
						Value:     "2",
					},
					{
						Timestamp: float64(start.Add(30*time.Minute).UnixMilli() / 1000.0),
						Value:     "1",
					},
					{
						Timestamp: float64(start.Add(40*time.Minute).UnixMilli() / 1000.0),
						Value:     "6",
					},
					{
						Timestamp: float64(start.Add(50*time.Minute).UnixMilli() / 1000.0),
						Value:     "4",
					},
					{
						Timestamp: float64(start.Add(60*time.Minute).UnixMilli() / 1000.0),
						Value:     "2",
					},
				},
			},
		}, nil, nil)

	out, _, err := fx.Run(cmd, []string{"123", "--type", "open_connections", "--start", "2022-11-01T00:00:00Z", "--end", "2022-11-01T01:00:00Z"})

	expOut := `Load Balancer:  	 Metric: open_connections 	 Start: 0001-01-01 00:00:00 +0000 UTC 	 End: 0001-01-01 00:00:00 +0000 UTC
 6.00 ┤                                                                 ╭─╮
 5.75 ┤                                                                ╭╯ ╰─╮
 5.50 ┤                                                               ╭╯    ╰─╮
 5.25 ┤                                                               │       ╰─╮
 5.01 ┤                                                              ╭╯         ╰─╮
 4.76 ┤                                                             ╭╯            ╰─╮
 4.51 ┤                                                            ╭╯               ╰─╮
 4.26 ┤                                                           ╭╯                  ╰─╮
 4.01 ┤               ╭─╮                                        ╭╯                     ╰─╮
 3.76 ┤             ╭─╯ ╰─╮                                      │                        ╰─╮
 3.52 ┤           ╭─╯     ╰─╮                                   ╭╯                          ╰─╮
 3.27 ┤         ╭─╯         ╰─╮                                ╭╯                             ╰─╮
 3.02 ┤       ╭─╯             ╰─╮                             ╭╯                                ╰─╮
 2.77 ┤     ╭─╯                 ╰─╮                          ╭╯                                   ╰─╮
 2.52 ┤  ╭──╯                     ╰──╮                      ╭╯                                      ╰──╮
 2.27 ┤╭─╯                           ╰─╮                    │                                          ╰─╮
 2.02 ┼╯                               ╰──╮                ╭╯                                            ╰
 1.78 ┤                                   ╰───╮           ╭╯
 1.53 ┤                                       ╰───╮      ╭╯
 1.28 ┤                                           ╰───╮ ╭╯
 1.03 ┤                                               ╰─╯


`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
