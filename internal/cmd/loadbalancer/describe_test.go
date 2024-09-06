package loadbalancer_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/loadbalancer"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := loadbalancer.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	lb := &hcloud.LoadBalancer{
		ID:   123,
		Name: "test",
		LoadBalancerType: &hcloud.LoadBalancerType{
			ID:                      123,
			Name:                    "lb11",
			Description:             "LB11",
			MaxServices:             5,
			MaxConnections:          10000,
			MaxTargets:              25,
			MaxAssignedCertificates: 10,
		},
		Created: time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
		PublicNet: hcloud.LoadBalancerPublicNet{
			Enabled: true,
			IPv4: hcloud.LoadBalancerPublicNetIPv4{
				IP: net.ParseIP("192.168.2.1"),
			},
			IPv6: hcloud.LoadBalancerPublicNetIPv6{
				IP: net.IPv6loopback,
			},
		},
		Algorithm: hcloud.LoadBalancerAlgorithm{
			Type: hcloud.LoadBalancerAlgorithmTypeLeastConnections,
		},
		IncludedTraffic: 20 * util.Tebibyte,
		IngoingTraffic:  10 * util.Tebibyte,
		OutgoingTraffic: 10 * util.Tebibyte,
	}

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(lb, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := fmt.Sprintf(`ID:				123
Name:				test
Created:			%s (%s)
Public Net:
  Enabled:			yes
  IPv4:				192.168.2.1
  IPv4 DNS PTR:			
  IPv6:				::1
  IPv6 DNS PTR:			
Private Net:
    No Private Network
Algorithm:			least_connections
Load Balancer Type:		lb11 (ID: 123)
  ID:				123
  Name:				lb11
  Description:			LB11
  Max Services:			5
  Max Connections:		10000
  Max Targets:			25
  Max assigned Certificates:	10
Services:
  No services
Targets:
  No targets
Traffic:
  Outgoing:	10 TiB
  Ingoing:	10 TiB
  Included:	20 TiB
Protection:
  Delete:	no
Labels:
  No labels
`, util.Datetime(lb.Created), humanize.Time(lb.Created))

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
