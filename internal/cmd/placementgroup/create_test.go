package placementgroup_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/placementgroup"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	opts := hcloud.PlacementGroupCreateOpts{
		Name:   "my Placement Group",
		Labels: map[string]string{},
		Type:   hcloud.PlacementGroupTypeSpread,
	}

	placementGroup := hcloud.PlacementGroup{
		ID:      897,
		Name:    opts.Name,
		Created: time.Now(),
		Labels:  opts.Labels,
		Type:    opts.Type,
	}

	fx.Client.PlacementGroupClient.EXPECT().
		Create(gomock.Any(), opts).
		Return(hcloud.PlacementGroupCreateResult{PlacementGroup: &placementGroup, Action: &hcloud.Action{ID: 321}}, nil, nil)

	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 321})

	out, _, err := fx.Run(cmd, []string{"--name", placementGroup.Name, "--type", string(placementGroup.Type)})

	expOut := `Placement group 897 created
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
