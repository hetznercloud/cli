package firewall_test

import (
	_ "embed"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/firewall"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/create_response.json
var createResponseJson string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := firewall.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.FirewallClient.EXPECT().
		Create(gomock.Any(), hcloud.FirewallCreateOpts{
			Name:   "test",
			Labels: make(map[string]string),
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
		Return(hcloud.FirewallCreateResult{
			Firewall: &hcloud.Firewall{
				ID:   123,
				Name: "test",
			},
			Actions: []*hcloud.Action{{ID: 321}},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}})

	out, _, err := fx.Run(cmd, []string{"--name", "test", "--rules-file", "testdata/rules.json"})

	expOut := "Firewall 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := firewall.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.FirewallClient.EXPECT().
		Create(gomock.Any(), hcloud.FirewallCreateOpts{
			Name:   "test",
			Labels: make(map[string]string),
		}).
		Return(hcloud.FirewallCreateResult{
			Firewall: &hcloud.Firewall{
				ID:      123,
				Name:    "test",
				Created: time.Date(2016, 1, 30, 23, 50, 0, 0, time.UTC),
				AppliedTo: []hcloud.FirewallResource{
					{Type: "server", Server: &hcloud.FirewallResourceServer{
						ID: 1,
					}},
				},
				Labels: make(map[string]string),
				Rules: []hcloud.FirewallRule{
					{
						Direction: "in",
						SourceIPs: []net.IPNet{},
						Protocol:  "tcp",
						Port:      hcloud.Ptr("22"),
					},
				},
			},
			Actions: []*hcloud.Action{{ID: 321}},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}})

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name", "test"})

	expOut := "Firewall 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJson, jsonOut)
}
