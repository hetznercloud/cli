package server_test

import (
	"crypto/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestIPv4(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.IPCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	buf := make([]byte, 4)
	_, _ = rand.Read(buf)
	ip := net.IP(buf)

	srv := &hcloud.Server{
		ID:   123,
		Name: "my-server",
		PublicNet: hcloud.ServerPublicNet{
			IPv4: hcloud.ServerPublicNetIPv4{
				IP: ip,
			},
		},
	}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)

	args := []string{"my-server"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, ip.String()+"\n", out)
}

func TestIPv6(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.IPCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	buf := make([]byte, 16)
	_, _ = rand.Read(buf[:8])
	ip := net.IP(buf)

	srv := &hcloud.Server{
		ID:   123,
		Name: "my-server",
		PublicNet: hcloud.ServerPublicNet{
			IPv6: hcloud.ServerPublicNetIPv6{
				IP: ip,
			},
		},
	}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)

	args := []string{"my-server", "-6"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, ip.String()+"1\n", out)
}
