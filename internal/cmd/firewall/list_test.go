package firewall_test

import (
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

	out, _, err := fx.Run(cmd, []string{})

	expOut := `ID    NAME   RULES COUNT   APPLIED TO COUNT
123   test   5 Rules       2 Servers | 0 Label Selectors
`

	assert.NoError(t, err)
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
				ID:        123,
				Name:      "test",
				Rules:     make([]hcloud.FirewallRule, 5),
				AppliedTo: make([]hcloud.FirewallResource, 2),
				Labels:    make(map[string]string),
			},
		}, nil)

	out, _, err := fx.Run(cmd, []string{"-o=json"})

	assert.NoError(t, err)
	assert.JSONEq(t, out, `
[
  {
    "id": 123,
    "name": "test",
    "labels": {},
    "created": "0001-01-01T00:00:00Z",
    "rules": [
      {
        "direction": "",
        "protocol": ""
      },
      {
        "direction": "",
        "protocol": ""
      },
      {
        "direction": "",
        "protocol": ""
      },
      {
        "direction": "",
        "protocol": ""
      },
      {
        "direction": "",
        "protocol": ""
      }
    ],
    "applied_to": [
      {
        "type": ""
      },
      {
        "type": ""
      }
    ]
  }
]`)
}
