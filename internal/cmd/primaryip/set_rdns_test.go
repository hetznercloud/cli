package primaryip_test

import (
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/primaryip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestSetRDNS(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := primaryip.SetRDNSCmd.CobraCommand(fx.State())
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
	fx.Client.RDNSClient.EXPECT().
		ChangeDNSPtr(
			gomock.Any(),
			&hcloud.PrimaryIP{ID: 13},
			net.ParseIP("192.168.0.1"),
			hcloud.Ptr("server.your-host.de"),
		).
		Return(
			action,
			&hcloud.Response{},
			nil,
		)

	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), action).Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"--hostname=server.your-host.de", "--ip=192.168.0.1", "13"})

	expOut := "Reverse DNS of Primary IP 13 changed\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
