package server_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestMetrics(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.MetricsCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	start := time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	srv := &hcloud.Server{ID: 123, Name: "my-server"}

	fx.Client.Server.EXPECT().
		Get(gomock.Any(), "123").
		Return(srv, nil, nil)
	fx.Client.Server.EXPECT().
		GetMetrics(gomock.Any(), srv, hcloud.ServerGetMetricsOpts{
			Start: start,
			End:   end,
			Types: []hcloud.ServerMetricType{hcloud.ServerMetricCPU},
		}).
		Return(&hcloud.ServerMetrics{
			TimeSeries: map[string][]hcloud.ServerMetricsValue{
				"cpu": {
					{
						Timestamp: float64(start.Add(0*time.Minute).UnixMilli() / 1000.0),
						Value:     "0.1",
					},
					{
						Timestamp: float64(start.Add(10*time.Minute).UnixMilli() / 1000.0),
						Value:     "1.2",
					},
					{
						Timestamp: float64(start.Add(20*time.Minute).UnixMilli() / 1000.0),
						Value:     "0.3",
					},
					{
						Timestamp: float64(start.Add(30*time.Minute).UnixMilli() / 1000.0),
						Value:     "0.8",
					},
					{
						Timestamp: float64(start.Add(40*time.Minute).UnixMilli() / 1000.0),
						Value:     "0.9",
					},
					{
						Timestamp: float64(start.Add(50*time.Minute).UnixMilli() / 1000.0),
						Value:     "0.1",
					},
					{
						Timestamp: float64(start.Add(60*time.Minute).UnixMilli() / 1000.0),
						Value:     "0.2",
					},
				},
			},
			Start: start,
			End:   end,
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"123", "--type", "cpu", "--start", "2022-11-01T00:00:00Z", "--end", "2022-11-01T01:00:00Z"})

	expOut := `Server: my-server 	 Metric: cpu 	 Start: 2022-11-01 00:00:00 +0000 UTC 	 End: 2022-11-01 01:00:00 +0000 UTC
 1.17 ┤               ╭─╮
 1.12 ┤              ╭╯ ╰╮
 1.07 ┤              │   ╰╮
 1.01 ┤             ╭╯    ╰╮
 0.96 ┤            ╭╯      ╰╮
 0.90 ┤           ╭╯        ╰╮                                        ╭──╮
 0.85 ┤          ╭╯          ╰╮                              ╭────────╯  ╰╮
 0.80 ┤          │            ╰╮                       ╭─────╯            ╰╮
 0.74 ┤         ╭╯             ╰╮                    ╭─╯                   ╰╮
 0.69 ┤        ╭╯               ╰╮                  ╭╯                      ╰╮
 0.64 ┤       ╭╯                 ╰╮               ╭─╯                        ╰╮
 0.58 ┤      ╭╯                   ╰╮            ╭─╯                           ╰╮
 0.53 ┤      │                     ╰╮         ╭─╯                              ╰─╮
 0.48 ┤     ╭╯                      ╰╮       ╭╯                                  ╰╮
 0.42 ┤    ╭╯                        ╰╮    ╭─╯                                    ╰╮
 0.37 ┤   ╭╯                          ╰╮ ╭─╯                                       ╰╮
 0.31 ┤  ╭╯                            ╰─╯                                          ╰╮
 0.26 ┤  │                                                                           ╰╮
 0.21 ┤ ╭╯                                                                            ╰╮               ╭──
 0.15 ┤╭╯                                                                              ╰╮      ╭───────╯
 0.10 ┼╯                                                                                ╰──────╯


`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
