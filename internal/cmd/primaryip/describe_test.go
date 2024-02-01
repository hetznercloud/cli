package primaryip_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

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
		Datacenter:   &hcloud.Datacenter{ID: 0, Location: &hcloud.Location{ID: 0}},
	}

	fx.Client.PrimaryIPClient.EXPECT().
		Get(gomock.Any(), "10").
		Return(primaryIP, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"10"})

	expOut := fmt.Sprintf(`ID:		10
Name:		test-net
Created:	%s (%s)
Type:		ipv4
IP:		192.168.0.1
Blocked:	yes
Auto delete:	no
Assignee:
  Not assigned
DNS:
  No reverse DNS entries
Protection:
  Delete:	no
Labels:
  No labels
Datacenter:
  ID:		0
  Name:		
  Description:	
  Location:
    Name:		
    Description:	
    Country:		
    City:		
    Latitude:		0.000000
    Longitude:		0.000000
`, util.Datetime(primaryIP.Created), humanize.Time(primaryIP.Created))

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
