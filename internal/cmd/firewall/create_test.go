package firewall

import (
	_ "embed"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/create_response.json
var createResponseJson string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.FirewallClient.EXPECT().
		Create(gomock.Any(), hcloud.FirewallCreateOpts{
			Name:   "test",
			Labels: make(map[string]string),
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

	out, _, err := fx.Run(cmd, []string{"--name", "test"})

	expOut := "Firewall 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(fx.State())
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
