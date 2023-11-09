package network

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/util"
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

	network := &hcloud.Network{
		ID:         123,
		Name:       "test",
		Created:    time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
		IPRange:    &net.IPNet{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(24, 32)},
		Protection: hcloud.NetworkProtection{Delete: true},
		Labels:     map[string]string{"key": "value"},
	}

	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(network, nil, nil)

	out, _, err := fx.Run(cmd, []string{"test"})

	expOut := fmt.Sprintf(`ID:		123
Name:		test
Created:	%s (%s)
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
`, util.Datetime(network.Created), humanize.Time(network.Created))

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
