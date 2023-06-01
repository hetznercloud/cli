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

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.DeleteCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	placementGroup := hcloud.PlacementGroup{
		ID:      897,
		Name:    "my Placement Group",
		Created: time.Now(),
		Labels:  map[string]string{"key": "value"},
		Servers: []int{4711, 4712},
		Type:    hcloud.PlacementGroupTypeSpread,
	}

	fx.Client.PlacementGroupClient.EXPECT().
		Get(gomock.Any(), placementGroup.Name).
		Return(&placementGroup, nil, nil)
	fx.Client.PlacementGroupClient.EXPECT().
		Delete(gomock.Any(), &placementGroup)

	_, err := fx.Run(cmd, []string{placementGroup.Name})

	assert.NoError(t, err)
}
