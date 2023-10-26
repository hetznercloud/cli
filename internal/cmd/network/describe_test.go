package network

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := DescribeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.Network{
			ID:         123,
			Name:       "test",
			Created:    time.Date(1905, 10, 6, 12, 0, 0, 0, time.UTC),
			IPRange:    &net.IPNet{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(24, 32)},
			Protection: hcloud.NetworkProtection{Delete: true},
			Labels:     map[string]string{"key": "value"},
		}, nil, nil)

	out, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:		123
Name:		test
Created:	Fri Oct  6 12:00:00 UTC 1905 (a long while ago)
IP Range:	10.0.0.0/24
Expose Routes to vSwitch: no
Subnets:
  No subnets
Routes:
  No routes
Protection:
  Delete:	yes
Labels:
  key: value
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
