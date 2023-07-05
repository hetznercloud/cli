package primaryip

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := describeCmd.CobraCommand(context.Background(), fx.Client, fx.TokenEnsurer)
	created := time.Date(1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	ip := net.ParseIP("192.168.0.1")
	fx.ExpectEnsureToken()
	fx.Client.PrimaryIPClient.EXPECT().
		Get(
			gomock.Any(),
			"10",
		).
		Return(&hcloud.PrimaryIP{
			ID:           10,
			Name:         "test-net",
			Type:         "ipv4",
			Created:      created,
			IP:           ip,
			Blocked:      true,
			AutoDelete:   false,
			AssigneeType: "server",
			Datacenter:   &hcloud.Datacenter{ID: 0, Location: &hcloud.Location{ID: 0}},
		},
			&hcloud.Response{},
			nil)

	out, err := fx.Run(cmd, []string{"10"})

	expOut := `ID:		10
Name:		test-net
Created:	%s (a long while ago)
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
`
	expOutFmt := fmt.Sprintf(expOut, util.Datetime(created))
	assert.NoError(t, err)
	assert.Equal(t, expOutFmt, out)
}
