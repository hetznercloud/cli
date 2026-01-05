package server_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := server.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.ServerClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.ServerListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
				Status:   []hcloud.ServerStatus{hcloud.ServerStatusRunning},
			},
		).
		Return([]*hcloud.Server{
			{
				ID:       123,
				Name:     "test",
				Status:   hcloud.ServerStatusRunning,
				Location: &hcloud.Location{Name: "fsn1"},
				PublicNet: hcloud.ServerPublicNet{
					IPv4: hcloud.ServerPublicNetIPv4{
						IP: net.ParseIP("192.168.2.1"),
					},
				},
				Created: time.Now().Add(-20 * time.Second),
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{"--status", "running"})

	expOut := `ID    NAME   STATUS    IPV4          IPV6   PRIVATE NET   LOCATION   AGE
123   test   running   192.168.2.1   -      -             fsn1       20s
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
