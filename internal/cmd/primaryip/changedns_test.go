package primaryip

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestChangeDNS(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := ChangeDNSCmd.CobraCommand(context.Background(), fx.Client, fx.TokenEnsurer, fx.ActionWaiter)
	action := &hcloud.Action{ID: 1}
	fx.ExpectEnsureToken()
	fx.Client.PrimaryIPClient.EXPECT().
		Get(
			gomock.Any(),
			"13",
		).
		Return(
			&hcloud.PrimaryIP{ID: 13},
			&hcloud.Response{},
			nil,
		)
	fx.Client.PrimaryIPClient.EXPECT().
		ChangeDNSPtr(
			gomock.Any(),
			hcloud.PrimaryIPChangeDNSPtrOpts{
				ID:     13,
				DNSPtr: "server.your-host.de",
				IP:     "192.168.0.1",
			},
		).
		Return(
			action,
			&hcloud.Response{},
			nil,
		)

	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), action).Return(nil)

	out, _, err := fx.Run(cmd, []string{"--hostname=server.your-host.de", "--ip=192.168.0.1", "13"})

	expOut := "Primary IP 13 DNS pointer: server.your-host.de associated to 192.168.0.1\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
