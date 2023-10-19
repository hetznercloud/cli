package primaryip

import (
	"context"
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
				PrimaryIP: &hcloud.PrimaryIP{ID: 1},
			},
			&hcloud.Response{},
			nil,
		)

	out, err := fx.Run(cmd, []string{"--name=my-ip", "--type=ipv4", "--datacenter=fsn1-dc14"})

	expOut := "Primary IP 1 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
