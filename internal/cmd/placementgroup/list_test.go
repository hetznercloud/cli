package placementgroup_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/placementgroup"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.ListCmd.CobraCommand(fx.State())
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
				Servers: []int64{4711, 4712},
				Type:    hcloud.PlacementGroupTypeSpread,
				Created: time.Now().Add(-10 * time.Second),
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{"--selector", "foo=bar"})

	expOut := `ID    NAME                 SERVERS     TYPE     AGE
897   my Placement Group   2 servers   spread   10s
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
