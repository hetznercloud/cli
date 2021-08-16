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

func TestAddLabel(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.LabelCmds.AddCobraCommand(
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
		Labels: map[string]string{"key": "value", "foo": "bar"},
	}

	fx.Client.PlacementGroupClient.EXPECT().
		Get(gomock.Any(), placementGroup.Name).
		Return(&placementGroup, nil, nil)
	fx.Client.PlacementGroupClient.EXPECT().
		Update(gomock.Any(), &hcloud.PlacementGroup{ID: placementGroup.ID}, opts).
		Return(&placementGroup, nil, nil)

	_, err := fx.Run(cmd, []string{placementGroup.Name, "foo=bar"})

	assert.NoError(t, err)
}

func TestRemoveLabel(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.LabelCmds.RemoveCobraCommand(
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
		Labels: map[string]string{},
	}

	fx.Client.PlacementGroupClient.EXPECT().
		Get(gomock.Any(), placementGroup.Name).
		Return(&placementGroup, nil, nil)
	fx.Client.PlacementGroupClient.EXPECT().
		Update(gomock.Any(), &hcloud.PlacementGroup{ID: placementGroup.ID}, opts).
		Return(&placementGroup, nil, nil)

	_, err := fx.Run(cmd, []string{placementGroup.Name, "key"})

	assert.NoError(t, err)
}
