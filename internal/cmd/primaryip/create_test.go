package primaryip

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

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
