package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := ListCmd.CobraCommand(context.Background(), fx.Client, fx.TokenEnsurer)

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
				ID:         123,
				Name:       "test",
				Status:     hcloud.ServerStatusRunning,
				Datacenter: &hcloud.Datacenter{Name: "fsn1-dc14"},
				PublicNet: hcloud.ServerPublicNet{
					IPv4: hcloud.ServerPublicNetIPv4{
						IP: net.ParseIP("192.168.2.1"),
					},
				},
				Created: time.Now().Add(-20 * time.Second),
			},
		}, nil)

	out, _, err := fx.Run(cmd, []string{"--status", "running"})

	expOut := `ID    NAME   STATUS    IPV4          IPV6   PRIVATE NET   DATACENTER   AGE
123   test   running   192.168.2.1   -      -             fsn1-dc14    20s
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
