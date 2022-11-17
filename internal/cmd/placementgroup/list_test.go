package placementgroup_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/cli/internal/cmd/placementgroup"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.ListCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		nil)
	fx.ExpectEnsureToken()

	fx.Client.PlacementGroupClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.PlacementGroupListOpts{
				ListOpts: hcloud.ListOpts{
					PerPage:       50,
					LabelSelector: "foo=bar",
				},
				Sort: []string{"id:asc"},
			},
		).
		Return([]*hcloud.PlacementGroup{
			{
				ID:      897,
				Name:    "my Placement Group",
				Labels:  map[string]string{"key": "value"},
				Servers: []int{4711, 4712},
				Type:    hcloud.PlacementGroupTypeSpread,
				Created: time.Now().Add(-10 * time.Second),
			},
		}, nil)

	out, err := fx.Run(cmd, []string{"--selector", "foo=bar"})

	expOut := `ID    NAME                 SERVERS     TYPE     AGE
897   my Placement Group   2 servers   spread   10s
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
