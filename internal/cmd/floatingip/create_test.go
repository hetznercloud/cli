package floatingip_test

import (
	_ "embed"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/floatingip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/create_response.json
var createResponseJSON string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := floatingip.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.FloatingIPClient.EXPECT().
		Create(gomock.Any(), hcloud.FloatingIPCreateOpts{
			Name:         hcloud.Ptr("myFloatingIP"),
			Type:         hcloud.FloatingIPTypeIPv4,
			HomeLocation: &hcloud.Location{Name: "fsn1"},
			Labels:       make(map[string]string),
			Description:  hcloud.Ptr(""),
		}).
		Return(hcloud.FloatingIPCreateResult{
			FloatingIP: &hcloud.FloatingIP{
				ID:   123,
				Name: "myFloatingIP",
				IP:   net.ParseIP("192.168.2.1"),
				Type: hcloud.FloatingIPTypeIPv4,
			},
			Action: nil,
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"--name", "myFloatingIP", "--type", "ipv4", "--home-location", "fsn1"})

	expOut := `Floating IP 123 created
IPv4: 192.168.2.1
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := floatingip.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.FloatingIPClient.EXPECT().
		Create(gomock.Any(), hcloud.FloatingIPCreateOpts{
			Name:         hcloud.Ptr("myFloatingIP"),
			Type:         hcloud.FloatingIPTypeIPv4,
			HomeLocation: &hcloud.Location{Name: "fsn1"},
			Labels:       make(map[string]string),
			Description:  hcloud.Ptr(""),
		}).
		Return(hcloud.FloatingIPCreateResult{
			FloatingIP: &hcloud.FloatingIP{
				ID:     123,
				Name:   "myFloatingIP",
				IP:     net.ParseIP("127.0.0.1"),
				Type:   hcloud.FloatingIPTypeIPv4,
				Labels: map[string]string{},
				Server: &hcloud.Server{ID: 1},
			},
			Action: nil,
		}, nil, nil)

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name", "myFloatingIP", "--type", "ipv4", "--home-location", "fsn1"})

	expOut := "Floating IP 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJSON, jsonOut)
}

func TestCreateProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := floatingip.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	floatingIP := &hcloud.FloatingIP{
		ID:   123,
		Name: "myFloatingIP",
		IP:   net.ParseIP("192.168.2.1"),
		Type: hcloud.FloatingIPTypeIPv4,
	}

	fx.Client.FloatingIPClient.EXPECT().
		Create(gomock.Any(), hcloud.FloatingIPCreateOpts{
			Name:         hcloud.Ptr("myFloatingIP"),
			Type:         hcloud.FloatingIPTypeIPv4,
			HomeLocation: &hcloud.Location{Name: "fsn1"},
			Labels:       make(map[string]string),
			Description:  hcloud.Ptr(""),
		}).
		Return(hcloud.FloatingIPCreateResult{
			FloatingIP: floatingIP,
			Action: &hcloud.Action{
				ID: 321,
			},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321}).Return(nil)

	fx.Client.FloatingIPClient.EXPECT().
		ChangeProtection(gomock.Any(), floatingIP, hcloud.FloatingIPChangeProtectionOpts{
			Delete: hcloud.Ptr(true),
		}).
		Return(&hcloud.Action{ID: 333}, nil, nil)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 333}).Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"--name", "myFloatingIP", "--type", "ipv4", "--home-location", "fsn1", "--enable-protection", "delete"})

	expOut := `Floating IP 123 created
Resource protection enabled for floating IP 123
IPv4: 192.168.2.1
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
