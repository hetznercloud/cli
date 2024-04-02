package firewall_test

import (
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/firewall"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := firewall.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.FirewallClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.FirewallListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.Firewall{
			{
				ID:        123,
				Name:      "test",
				Rules:     make([]hcloud.FirewallRule, 5),
				AppliedTo: make([]hcloud.FirewallResource, 2),
				Labels:    make(map[string]string),
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID    NAME   RULES COUNT   APPLIED TO COUNT
123   test   5 Rules       2 Servers | 0 Label Selectors
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestListJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := firewall.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.FirewallClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.FirewallListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.Firewall{
			{
				ID:      42,
				Name:    "my-resource",
				Created: time.Date(2016, 1, 30, 23, 55, 0, 0, time.UTC),
				Labels: map[string]string{
					"environment":           "prod",
					"example.com/my":        "label",
					"just-a-key":            "",
					"just-a-key-with-value": "value",
				},
				Rules: []hcloud.FirewallRule{
					{
						Direction: "in",
						SourceIPs: []net.IPNet{
							{
								IP:   net.IP{28, 239, 13, 1},
								Mask: net.IPMask{255, 255, 255, 255},
							},
							{
								IP:   net.IP{28, 239, 14, 0},
								Mask: net.IPMask{255, 255, 255, 0},
							},
							{
								IP:   net.ParseIP("ff21:1eac:9a3b:ee58:5ca:990c:8bc9:c03b"),
								Mask: net.CIDRMask(128, 128),
							},
						},
						DestinationIPs: []net.IPNet{},
						Protocol:       hcloud.FirewallRuleProtocolTCP,
						Port:           hcloud.Ptr("80"),
					},
				},
				AppliedTo: []hcloud.FirewallResource{
					{
						Type: "server",
						Server: &hcloud.FirewallResourceServer{
							ID: 42,
						},
					},
				},
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{"-o=json"})

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.JSONEq(t, out, `
[
  {
    "id": 42,
    "name": "my-resource",
    "labels": {
      "environment": "prod",
      "example.com/my": "label",
      "just-a-key": "",
      "just-a-key-with-value": "value"
    },
    "created": "2016-01-30T23:55:00Z",
    "rules": [
      {
        "direction": "in",
        "source_ips": [
          "28.239.13.1/32",
          "28.239.14.0/24",
          "ff21:1eac:9a3b:ee58:5ca:990c:8bc9:c03b/128"
        ],
        "destination_ips": [],
        "protocol": "tcp",
        "port": "80",
        "description": null
      }
    ],
    "applied_to": [
      {
        "type": "server",
        "server": {
          "id": 42
        }
      }
    ]
  }
]`)
}
