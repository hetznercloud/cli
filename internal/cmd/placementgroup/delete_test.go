package placementgroup_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/placementgroup"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	placementGroup := hcloud.PlacementGroup{
		ID:      897,
		Name:    "my Placement Group",
		Created: time.Now(),
		Labels:  map[string]string{"key": "value"},
		Servers: []int64{4711, 4712},
		Type:    hcloud.PlacementGroupTypeSpread,
	}

	fx.Client.PlacementGroupClient.EXPECT().
		Get(gomock.Any(), placementGroup.Name).
		Return(&placementGroup, nil, nil)
	fx.Client.PlacementGroupClient.EXPECT().
		Delete(gomock.Any(), &placementGroup)

	out, errOut, err := fx.Run(cmd, []string{placementGroup.Name})

	expOut := "placement group my Placement Group deleted\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
