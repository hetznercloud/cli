package loadbalancer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/loadbalancer"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestAddService(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.AddServiceCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		AddService(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerAddServiceOpts{
			Protocol:        hcloud.LoadBalancerServiceProtocolHTTP,
			ListenPort:      new(80),
			DestinationPort: new(8080),
			HTTP: &hcloud.LoadBalancerAddServiceOptsHTTP{
				StickySessions: new(false),
				RedirectHTTP:   new(false),
				TimeoutIdle:    hcloud.Ptr(60 * time.Second),
			},
			Proxyprotocol: new(false),
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"123", "--protocol", "http", "--listen-port", "80", "--destination-port", "8080", "--http-timeout-idle", "60s"})

	expOut := "Service was added to Load Balancer 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestAddServiceWithHealthCheck(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.AddServiceCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.CertificateClient.EXPECT().
		Get(gomock.Any(), "1").
		Return(&hcloud.Certificate{ID: 1}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		AddService(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerAddServiceOpts{
			Protocol:        hcloud.LoadBalancerServiceProtocolHTTP,
			ListenPort:      new(80),
			DestinationPort: new(8080),
			HTTP: &hcloud.LoadBalancerAddServiceOptsHTTP{
				StickySessions: new(true),
				RedirectHTTP:   new(true),
				CookieName:     new("test"),
				Certificates:   []*hcloud.Certificate{{ID: 1}},
				CookieLifetime: hcloud.Ptr(10 * time.Minute),
			},
			Proxyprotocol: new(false),
			HealthCheck: &hcloud.LoadBalancerAddServiceOptsHealthCheck{
				Protocol: hcloud.LoadBalancerServiceProtocolHTTP,
				Port:     new(80),
				Interval: hcloud.Ptr(10 * time.Second),
				Timeout:  hcloud.Ptr(5 * time.Second),
				Retries:  new(2),
				HTTP: &hcloud.LoadBalancerAddServiceOptsHealthCheckHTTP{
					Domain:      new("example.com"),
					Path:        new("/health"),
					StatusCodes: []string{"200"},
					Response:    new("OK"),
					TLS:         new(true),
				},
			},
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{
		"123",
		"--protocol", "http",
		"--listen-port", "80",
		"--destination-port", "8080",
		"--http-redirect-http=true",
		"--http-sticky-sessions=true",
		"--http-cookie-name", "test",
		"--http-cookie-lifetime", "10m",
		"--http-certificates", "1",
		"--health-check-protocol", "http",
		"--health-check-port", "80",
		"--health-check-interval", "10s",
		"--health-check-timeout", "5s",
		"--health-check-retries", "2",
		"--health-check-http-domain", "example.com",
		"--health-check-http-path", "/health",
		"--health-check-http-status-codes", "200",
		"--health-check-http-response", "OK",
		"--health-check-http-tls=true",
	})

	expOut := "Service was added to Load Balancer 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
