package placementgroup_test

import (
	_ "embed"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/placementgroup"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/create_response.json
var createResponseJson string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.CreateCmd.CobraCommand(fx.State())
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
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321})

	out, errOut, err := fx.Run(cmd, []string{"--name", placementGroup.Name, "--type", string(placementGroup.Type)})

	expOut := `Placement group 897 created
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := placementgroup.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	opts := hcloud.PlacementGroupCreateOpts{
		Name:   "myPlacementGroup",
		Labels: map[string]string{},
		Type:   hcloud.PlacementGroupTypeSpread,
	}

	fx.Client.PlacementGroupClient.EXPECT().
		Create(gomock.Any(), opts).
		Return(hcloud.PlacementGroupCreateResult{
			PlacementGroup: &hcloud.PlacementGroup{
				ID:      897,
				Name:    "myPlacementGroup",
				Created: time.Date(2016, 1, 30, 23, 50, 0, 0, time.UTC),
				Servers: []int64{1, 2, 3},
				Labels:  make(map[string]string),
				Type:    hcloud.PlacementGroupTypeSpread,
			},
			Action: &hcloud.Action{ID: 321},
		}, nil, nil)

	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321})

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name", "myPlacementGroup", "--type", "spread"})

	expOut := "Placement group 897 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJson, jsonOut)
}
