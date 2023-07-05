package placementgroup_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/cli/internal/cmd/placementgroup"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestUpdateName(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.UpdateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	placementGroup := hcloud.PlacementGroup{
		ID:      897,
		Name:    "my Placement Group",
		Created: time.Now(),
		Labels:  map[string]string{"key": "value"},
		Servers: []int{4711, 4712},
		Type:    hcloud.PlacementGroupTypeSpread,
	}

	opts := hcloud.PlacementGroupUpdateOpts{
		Name: "new placement group name",
	}

	fx.Client.PlacementGroupClient.EXPECT().
		Get(gomock.Any(), placementGroup.Name).
		Return(&placementGroup, nil, nil)
	fx.Client.PlacementGroupClient.EXPECT().
		Update(gomock.Any(), &placementGroup, opts).
		Return(&placementGroup, nil, nil)

	_, err := fx.Run(cmd, []string{placementGroup.Name, "--name", opts.Name})

	assert.NoError(t, err)
}
