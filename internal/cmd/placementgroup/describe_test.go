package placementgroup_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/cli/internal/cmd/placementgroup"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.DescribeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
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

	out, err := fx.Run(cmd, []string{placementGroup.Name})

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
`, util.Datetime(placementGroup.Created),
		humanize.Time(placementGroup.Created),
	)

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
