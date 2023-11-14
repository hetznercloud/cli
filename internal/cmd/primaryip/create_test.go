package primaryip

import (
	"context"
	_ "embed"
	"net"
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

	cmd := CreateCmd.CobraCommand(context.Background(), fx.Client, fx.TokenEnsurer, fx.ActionWaiter)
	fx.ExpectEnsureToken()
	fx.Client.PrimaryIPClient.EXPECT().
		Create(
			gomock.Any(),
			hcloud.PrimaryIPCreateOpts{
				Name:         "my-ip",
				Type:         "ipv4",
				Datacenter:   "fsn1-dc14",
				AssigneeType: "server",
			},
		).
		Return(
			&hcloud.PrimaryIPCreateResult{
				PrimaryIP: &hcloud.PrimaryIP{
					ID:   1,
					IP:   net.ParseIP("192.168.2.1"),
					Type: hcloud.PrimaryIPTypeIPv4,
				},
				Action: &hcloud.Action{ID: 321},
			},
			&hcloud.Response{},
			nil,
		)

	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 321})

	out, _, err := fx.Run(cmd, []string{"--name=my-ip", "--type=ipv4", "--datacenter=fsn1-dc14"})

	expOut := `Primary IP 1 created
IPv4: 192.168.2.1
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	response, err := testutil.MockResponse(&schema.PrimaryIPCreateResponse{
		PrimaryIP: schema.PrimaryIP{
			ID:   1,
			Name: "my-ip",
			IP:   "192.168.2.1",
			Type: "ipv4",
			Datacenter: schema.Datacenter{
				ID:       1,
				Name:     "fsn1-dc14",
				Location: schema.Location{ID: 1, Name: "fsn1"},
			},
			Created:      time.Date(2016, 1, 30, 23, 50, 0, 0, time.UTC),
			Labels:       make(map[string]string),
			AutoDelete:   true,
			AssigneeID:   1,
			AssigneeType: "server",
			DNSPtr:       make([]schema.PrimaryIPDNSPTR, 0),
		},
		Action: &schema.Action{
			ID: 321,
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	fx.Client.PrimaryIPClient.EXPECT().
		Create(
			gomock.Any(),
			hcloud.PrimaryIPCreateOpts{
				Name:         "my-ip",
				Type:         "ipv4",
				Datacenter:   "fsn1-dc14",
				AssigneeType: "server",
			},
		).
		Return(
			&hcloud.PrimaryIPCreateResult{
				PrimaryIP: &hcloud.PrimaryIP{
					ID:   1,
					IP:   net.ParseIP("192.168.2.1"),
					Type: hcloud.PrimaryIPTypeIPv4,
				},
				Action: &hcloud.Action{ID: 321},
			}, response, nil)

	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 321})

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name=my-ip", "--type=ipv4", "--datacenter=fsn1-dc14"})

	expOut := "Primary IP 1 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJson, jsonOut)
}
