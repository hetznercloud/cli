package loadbalancer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateService(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := UpdateServiceCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		UpdateService(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, 80, hcloud.LoadBalancerUpdateServiceOpts{
			DestinationPort: hcloud.Ptr(8080),
			Protocol:        hcloud.LoadBalancerServiceProtocolTCP,
			Proxyprotocol:   hcloud.Ptr(true),
			HTTP: &hcloud.LoadBalancerUpdateServiceOptsHTTP{
				RedirectHTTP:   hcloud.Ptr(true),
				StickySessions: hcloud.Ptr(true),
				CookieName:     hcloud.Ptr("test"),
				CookieLifetime: hcloud.Ptr(10 * time.Minute),
				Certificates:   []*hcloud.Certificate{{ID: 1}},
			},
			HealthCheck: &hcloud.LoadBalancerUpdateServiceOptsHealthCheck{
				Protocol: hcloud.LoadBalancerServiceProtocolTCP,
				Port:     hcloud.Ptr(8080),
				Interval: hcloud.Ptr(10 * time.Second),
				Timeout:  hcloud.Ptr(5 * time.Second),
				Retries:  hcloud.Ptr(2),
				HTTP: &hcloud.LoadBalancerUpdateServiceOptsHealthCheckHTTP{
					Domain:      hcloud.Ptr("example.com"),
					Path:        hcloud.Ptr("/health"),
					StatusCodes: []string{"200"},
					Response:    hcloud.Ptr("OK"),
					TLS:         hcloud.Ptr(true),
				},
			},
		}).
		Return(&hcloud.Action{ID: 321}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{
		"123",
		"--listen-port", "80",
		"--destination-port", "8080",
		"--protocol", "tcp",
		"--proxy-protocol=true",
		"--http-redirect-http=true",
		"--http-sticky-sessions=true",
		"--http-cookie-name", "test",
		"--http-cookie-lifetime", "10m",
		"--http-certificates", "1",
		"--health-check-protocol", "tcp",
		"--health-check-port", "8080",
		"--health-check-interval", "10s",
		"--health-check-timeout", "5s",
		"--health-check-retries", "2",
		"--health-check-http-domain", "example.com",
		"--health-check-http-path", "/health",
		"--health-check-http-status-codes", "200",
		"--health-check-http-response", "OK",
		"--health-check-http-tls=true",
	})

	expOut := "Service 80 on Load Balancer 123 was updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
