package primaryip_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/primaryip"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := primaryip.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	primaryIP := &hcloud.PrimaryIP{
		ID:           10,
		Name:         "test-net",
		Type:         "ipv4",
		Created:      time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
		IP:           net.ParseIP("192.168.0.1"),
		Blocked:      true,
		AutoDelete:   false,
		AssigneeType: "server",
		Location: &hcloud.Location{
			ID:          3,
			Name:        "hel1",
			Description: "Helsinki DC Park 1",
			NetworkZone: "eu-central",
			Country:     "FI",
			City:        "Helsinki",
			Latitude:    60.169855,
			Longitude:   24.938379,
		},
	}

	fx.Client.PrimaryIPClient.EXPECT().
		Get(gomock.Any(), "10").
		Return(primaryIP, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"10"})

	expOut := fmt.Sprintf(`ID:           10
Name:         test-net
Created:      %s (%s)
Type:         ipv4
IP:           192.168.0.1
Blocked:      yes
Auto delete:  no

Assignee:
  Not assigned

DNS:
  No reverse DNS entries

Protection:
  Delete:  no

Labels:
  No labels

Location:
  ID:            3
  Name:          hel1
  Description:   Helsinki DC Park 1
  Network Zone:  eu-central
  Country:       FI
  City:          Helsinki
  Latitude:      60.169855
  Longitude:     24.938379
`, util.Datetime(primaryIP.Created), humanize.Time(primaryIP.Created))

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
