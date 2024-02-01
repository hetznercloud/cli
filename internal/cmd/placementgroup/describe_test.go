package placementgroup_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/placementgroup"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := placementgroup.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	placementGroup := hcloud.PlacementGroup{
		ID:      897,
		Name:    "my Placement Group",
		Created: time.Date(2021, 07, 23, 10, 0, 0, 0, time.UTC),
		Labels:  map[string]string{"key": "value"},
		Servers: []int64{4711, 4712},
		Type:    hcloud.PlacementGroupTypeSpread,
	}

	fx.Client.PlacementGroupClient.EXPECT().
		Get(gomock.Any(), placementGroup.Name).
		Return(&placementGroup, nil, nil)
	fx.Client.ServerClient.EXPECT().
		ServerName(int64(4711)).
		Return("server1")
	fx.Client.ServerClient.EXPECT().
		ServerName(int64(4712)).
		Return("server2")

	out, errOut, err := fx.Run(cmd, []string{placementGroup.Name})

	expOut := fmt.Sprintf(`ID:		897
Name:		my Placement Group
Created:	%s (%s)
Labels:
  key: value
Servers:
  - Server ID:		4711
    Server Name:	server1
  - Server ID:		4712
    Server Name:	server2
Type:		spread
`, util.Datetime(placementGroup.Created), humanize.Time(placementGroup.Created))

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
