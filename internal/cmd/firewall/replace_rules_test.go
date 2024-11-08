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

func TestReplaceRules(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := firewall.ReplaceRulesCmd.CobraCommand(fx.State())
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
			Rules: []hcloud.FirewallRule{
				{
					Direction: hcloud.FirewallRuleDirectionIn,
					SourceIPs: []net.IPNet{
						{IP: net.IP{28, 239, 13, 1}, Mask: net.IPMask{255, 255, 255, 255}},
						{IP: net.IP{28, 239, 14, 0}, Mask: net.IPMask{255, 255, 255, 0}},
						{
							IP:   net.IP{0xff, 0x21, 0x1e, 0xac, 0x9a, 0x3b, 0xee, 0x58, 0x05, 0xca, 0x99, 0x0c, 0x8b, 0xc9, 0xc0, 0x3b},
							Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
						},
					},
					DestinationIPs: nil,
					Protocol:       hcloud.FirewallRuleProtocolTCP,
					Port:           hcloud.Ptr("80"),
					Description:    hcloud.Ptr("Allow port 80"),
				},
				{
					Direction: hcloud.FirewallRuleDirectionIn,
					SourceIPs: []net.IPNet{
						{IP: net.IP{0, 0, 0, 0}, Mask: net.IPMask{0, 0, 0, 0}},
						{
							IP:   net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
							Mask: net.IPMask{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						},
					},
					DestinationIPs: nil,
					Protocol:       hcloud.FirewallRuleProtocolTCP,
					Port:           hcloud.Ptr("443"),
					Description:    hcloud.Ptr("Allow port 443"),
				},
				{
					Direction: hcloud.FirewallRuleDirectionOut,
					SourceIPs: nil,
					DestinationIPs: []net.IPNet{
						{IP: net.IP{28, 239, 13, 1}, Mask: net.IPMask{255, 255, 255, 255}},
						{IP: net.IP{28, 239, 14, 0}, Mask: net.IPMask{255, 255, 255, 0}},
						{
							IP:   net.IP{0xff, 0x21, 0x1e, 0xac, 0x9a, 0x3b, 0xee, 0x58, 0x05, 0xca, 0x99, 0x0c, 0x8b, 0xc9, 0xc0, 0x3b},
							Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
						},
					},
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port:     hcloud.Ptr("80"),
				},
			},
		}).
		Return([]*hcloud.Action{{ID: 123}, {ID: 321}}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 123}, {ID: 321}}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"--rules-file", "testdata/rules.json", "test"})

	expOut := "Firewall Rules for Firewall 123 updated\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
