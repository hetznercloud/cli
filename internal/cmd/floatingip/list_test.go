package floatingip_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/floatingip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := floatingip.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.FloatingIPClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.FloatingIPListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.FloatingIP{
			{
				ID:           123,
				Type:         hcloud.FloatingIPTypeIPv4,
				Name:         "test",
				HomeLocation: &hcloud.Location{Name: "fsn1"},
				IP:           net.ParseIP("192.168.2.1"),
				Created:      time.Now().Add(-10 * time.Minute),
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID    TYPE   NAME   DESCRIPTION   IP            HOME   SERVER   DNS   AGE
123   ipv4   test   -             192.168.2.1   fsn1   -        -     10m
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
