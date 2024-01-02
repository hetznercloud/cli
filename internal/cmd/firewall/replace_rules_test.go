package firewall

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestReplaceRules(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := ReplaceRulesCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	firewall := &hcloud.Firewall{
		ID:   123,
		Name: "test",
	}

	fx.Client.FirewallClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(firewall, nil, nil)
	fx.Client.FirewallClient.EXPECT().
		SetRules(gomock.Any(), firewall, hcloud.FirewallSetRulesOpts{
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
			},
		}).
		Return([]*hcloud.Action{{ID: 123}, {ID: 321}}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 123}, {ID: 321}}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"--rules-file", "testdata/rules.json", "test"})

	expOut := "Firewall Rules for Firewall 123 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
