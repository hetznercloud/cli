package floatingip_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/floatingip"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := floatingip.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	floatingIP := &hcloud.FloatingIP{
		ID:           123,
		Type:         hcloud.FloatingIPTypeIPv4,
		Name:         "test",
		Server:       &hcloud.Server{ID: 321},
		HomeLocation: &hcloud.Location{Name: "fsn1"},
		IP:           net.ParseIP("192.168.2.1"),
		Labels: map[string]string{
			"key": "value",
		},
		Created: time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
	}

	fx.Client.FloatingIPClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(floatingIP, nil, nil)
	fx.Client.ServerClient.EXPECT().
		ServerName(int64(321)).
		Return("myServer")

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := fmt.Sprintf(`ID:		123
Type:		ipv4
Name:		test
Description:	-
Created:	%s (%s)
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
`, util.Datetime(floatingIP.Created), humanize.Time(floatingIP.Created))

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
