package firewall_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/firewall"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestAddRule(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := firewall.AddRuleCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fw := &hcloud.Firewall{
		ID:   123,
		Name: "test",
	}

	fx.Client.FirewallClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(fw, nil, nil)
	fx.Client.FirewallClient.EXPECT().
		SetRules(gomock.Any(), fw, hcloud.FirewallSetRulesOpts{
			Rules: []hcloud.FirewallRule{{
				Direction:      hcloud.FirewallRuleDirectionIn,
				SourceIPs:      []net.IPNet{{IP: net.IP{0, 0, 0, 0}, Mask: net.IPMask{0, 0, 0, 0}}, {IP: net.IP{127, 0, 0, 1}, Mask: net.IPMask{255, 255, 255, 255}}},
				DestinationIPs: []net.IPNet{},
				Protocol:       hcloud.FirewallRuleProtocolTCP,
				Port:           hcloud.Ptr("80"),
				Description:    hcloud.Ptr("http"),
			}},
		}).
		Return([]*hcloud.Action{{ID: 123}, {ID: 321}}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 123}, {ID: 321}}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"--direction", "in", "--protocol", "tcp", "--source-ips", "0.0.0.0/0,127.0.0.1/32", "--port", "80", "--description", "http", "test"})

	expOut := "Firewall Rules for Firewall 123 updated\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
