package firewall_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/firewall"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := firewall.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fw := &hcloud.Firewall{
		ID:   123,
		Name: "test",
		Rules: []hcloud.FirewallRule{
			{
				Direction:   hcloud.FirewallRuleDirectionIn,
				Description: hcloud.Ptr("ssh"),
				Port:        hcloud.Ptr("22"),
				Protocol:    hcloud.FirewallRuleProtocolTCP,
			},
		},
		AppliedTo: []hcloud.FirewallResource{
			{
				Type: hcloud.FirewallResourceTypeServer,
				Server: &hcloud.FirewallResourceServer{
					ID: 321,
				},
			},
		},
		Labels: map[string]string{
			"key": "value",
		},
		Created: time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
	}

	fx.Client.FirewallClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(fw, nil, nil)
	fx.Client.ServerClient.EXPECT().
		ServerName(int64(321)).
		Return("myServer")

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := fmt.Sprintf(`ID:		123
Name:		test
Created:	%s (%s)
Labels:
  key: value
Rules:
  - Direction:		in
    Description:	ssh
    Protocol:		tcp
    Port:		22
    Origin IPs:
Applied To:
  - Type:		server
    Server ID:		321
    Server Name:	myServer
`, util.Datetime(fw.Created), humanize.Time(fw.Created))

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
