package network_test

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/cli/internal/cmd/network"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := network.ListCmd.CobraCommand(context.Background(), fx.Client, fx.TokenEnsurer)

	fx.ExpectEnsureToken()
	fx.Client.NetworkClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.NetworkListOpts{
				ListOpts: hcloud.ListOpts{
					PerPage:       50,
					LabelSelector: "foo=bar",
				},
			},
		).
		Return([]*hcloud.Network{
			{
				ID:      123,
				Name:    "test-net",
				IPRange: &net.IPNet{IP: net.ParseIP("192.0.2.1"), Mask: net.CIDRMask(24, 32)},
				Servers: []*hcloud.Server{{ID: 3421}},
			},
		},
			nil)

	out, err := fx.Run(cmd, []string{"--selector", "foo=bar"})

	expOut := `ID    NAME       IP RANGE       SERVERS
123   test-net   192.0.2.1/24   1 server
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
