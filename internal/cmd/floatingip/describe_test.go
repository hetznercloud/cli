package floatingip

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

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := DescribeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.FloatingIPClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.FloatingIP{
			ID:           123,
			Type:         hcloud.FloatingIPTypeIPv4,
			Name:         "test",
			Server:       &hcloud.Server{ID: 321},
			HomeLocation: &hcloud.Location{Name: "fsn1"},
			IP:           net.ParseIP("192.168.2.1"),
			Labels: map[string]string{
				"key": "value",
			},
			Created: time.Date(2036, 8, 20, 12, 0, 0, 0, time.UTC),
		}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		ServerName(int64(321)).
		Return("myServer")

	out, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:		123
Type:		ipv4
Name:		test
Description:	-
Created:	Wed Aug 20 12:00:00 UTC 2036 (12 years from now)
IP:		192.168.2.1
Blocked:	no
Home Location:	fsn1
Server:
  ID:	321
  Name:	myServer
DNS:
  No reverse DNS entries
Protection:
  Delete:	no
Labels:
  key: value
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
