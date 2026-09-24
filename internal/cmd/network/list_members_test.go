package network_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/network"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestListMembers(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := network.ListMembersCmd.CobraCommand(fx.State())

	networkRes := &hcloud.Network{
		ID:      123,
		Name:    "test-network",
		IPRange: &net.IPNet{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(16, 32)},
	}

	members := []*hcloud.NetworkMember{
		{
			Type:     hcloud.NetworkMemberTypeServer,
			ID:       42,
			IP:       net.IPv4(10, 0, 0, 2),
			Status:   hcloud.NetworkMemberStatusOK,
			AliasIPs: []net.IP{net.IPv4(10, 0, 0, 3), net.IPv4(10, 0, 0, 4)},
			Subnet: &net.IPNet{
				IP:   net.IPv4(10, 0, 0, 0),
				Mask: net.CIDRMask(24, 32),
			},
		},
		{
			Type:   hcloud.NetworkMemberTypeLoadBalancer,
			ID:     7,
			IP:     net.IPv4(10, 0, 0, 5),
			Status: hcloud.NetworkMemberStatusAttaching,
			Subnet: &net.IPNet{
				IP:   net.IPv4(10, 0, 0, 0),
				Mask: net.CIDRMask(24, 32),
			},
		},
	}

	fx.ExpectEnsureToken()
	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "test-network").
		Return(networkRes, nil, nil)
	fx.Client.NetworkClient.EXPECT().
		AllMembersWithOpts(
			gomock.Any(),
			networkRes,
			hcloud.NetworkMemberListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return(members, nil)

	out, errOut, err := fx.Run(cmd, []string{"test-network"})

	expOut := `TYPE            ID   IP         STATUS      ALIAS IPS   SUBNET
server          42   10.0.0.2   ok          10.0.0.3    10.0.0.0/24
                                            10.0.0.4
load_balancer   7    10.0.0.5   attaching   -           10.0.0.0/24
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
