package network_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/network"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := network.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.NetworkClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.NetworkListOpts{
				ListOpts: hcloud.ListOpts{
					PerPage:       50,
					LabelSelector: "foo=bar",
				},
				Sort: nil, // Networks do not support sorting
			},
		).
		Return([]*hcloud.Network{
			{
				ID:      123,
				Name:    "test-net",
				IPRange: &net.IPNet{IP: net.ParseIP("192.0.2.1"), Mask: net.CIDRMask(24, 32)},
				Servers: []*hcloud.Server{{ID: 3421}},
				Created: time.Now().Add(-10 * time.Second),
			},
		},
			nil)

	out, errOut, err := fx.Run(cmd, []string{"--selector", "foo=bar"})

	expOut := `ID    NAME       IP RANGE       SERVERS    AGE
123   test-net   192.0.2.1/24   1 server   10s
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
