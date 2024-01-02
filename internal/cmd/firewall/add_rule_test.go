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

func TestAddRule(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := AddRuleCmd.CobraCommand(
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
			Rules: []hcloud.FirewallRule{{
				Direction:      hcloud.FirewallRuleDirectionIn,
				SourceIPs:      []net.IPNet{{IP: net.IP{0, 0, 0, 0}, Mask: net.IPMask{0, 0, 0, 0}}},
				DestinationIPs: nil,
				Protocol:       hcloud.FirewallRuleProtocolTCP,
				Port:           hcloud.Ptr("80"),
				Description:    hcloud.Ptr("http"),
			}},
		}).
		Return([]*hcloud.Action{{ID: 123}, {ID: 321}}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 123}, {ID: 321}}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"--direction", "in", "--protocol", "tcp", "--source-ips", "0.0.0.0/0", "--port", "80", "--description", "http", "test"})

	expOut := "Firewall Rules for Firewall 123 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
