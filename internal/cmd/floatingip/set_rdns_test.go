package floatingip_test

import (
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/floatingip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestSetRDNS(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := floatingip.SetRDNSCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	floatingIP := &hcloud.FloatingIP{
		ID: 123,
		IP: net.ParseIP("192.168.2.1"),
	}

	fx.Client.FloatingIPClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(floatingIP, nil, nil)
	fx.Client.RDNSClient.EXPECT().
		ChangeDNSPtr(gomock.Any(), floatingIP, floatingIP.IP, hcloud.Ptr("example.com")).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"--hostname", "example.com", "test"})

	expOut := "Reverse DNS of Floating IP test changed\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
