package server_test

import (
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestAttachToNetwork(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.AttachToNetworkCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{ID: 123, Name: "my-server"}
	n := &hcloud.Network{ID: 456, Name: "my-network"}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)
	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "my-network").
		Return(n, nil, nil)
	fx.Client.ServerClient.EXPECT().
		AttachToNetwork(gomock.Any(), srv, hcloud.ServerAttachToNetworkOpts{
			Network: n,
			IP:      net.ParseIP("192.168.0.1"),
			AliasIPs: []net.IP{
				net.ParseIP("10.0.1.2"),
				net.ParseIP("10.0.1.3"),
			},
		}).
		Return(&hcloud.Action{ID: 789}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"my-server", "--network", "my-network", "--ip", "192.168.0.1", "--alias-ips", "10.0.1.2,10.0.1.3"}
	out, errOut, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Server 123 attached to network 456\n", out)
}
