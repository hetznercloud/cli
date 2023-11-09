package firewall

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

//go:embed testdata/create_response.json
var createResponseJson string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
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
		WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 321}})

	out, _, err := fx.Run(cmd, []string{"--name", "test"})

	expOut := "Firewall 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	response, err := testutil.MockResponse(&schema.FirewallCreateResponse{
		Firewall: schema.Firewall{
			ID:      123,
			Name:    "test",
			Created: time.Date(2016, 1, 30, 23, 50, 0, 0, time.UTC),
			AppliedTo: []schema.FirewallResource{
				{Type: "server", Server: &schema.FirewallResourceServer{
					ID: 1,
				}},
			},
			Labels: make(map[string]string),
			Rules: []schema.FirewallRule{
				{
					Direction: "in",
					SourceIPs: make([]string, 0),
					Protocol:  "tcp",
					Port:      hcloud.Ptr("22"),
				},
			},
		},
		Actions: make([]schema.Action, 0),
	})

	if err != nil {
		t.Fatal(err)
	}

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
		}, response, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 321}})

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name", "test"})

	expOut := "Firewall 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJson, jsonOut)
}
