package placementgroup_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

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

	fx.Client.PlacementGroup.EXPECT().
		Get(gomock.Any(), placementGroup.Name).
		Return(&placementGroup, nil, nil)
	fx.Client.PlacementGroup.EXPECT().
		Delete(gomock.Any(), &placementGroup)

	out, errOut, err := fx.Run(cmd, []string{placementGroup.Name})

	expOut := "placement group my Placement Group deleted\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	groups := []*hcloud.PlacementGroup{
		{
			ID:   123,
			Name: "test1",
		},
		{
			ID:   456,
			Name: "test2",
		},
		{
			ID:   789,
			Name: "test3",
		},
	}

	expOutBuilder := strings.Builder{}

	var names []string
	for _, pg := range groups {
		names = append(names, pg.Name)
		expOutBuilder.WriteString(fmt.Sprintf("placement group %s deleted\n", pg.Name))
		fx.Client.PlacementGroup.EXPECT().
			Get(gomock.Any(), pg.Name).
			Return(pg, nil, nil)
		fx.Client.PlacementGroup.EXPECT().
			Delete(gomock.Any(), pg).
			Return(nil, nil)
	}

	out, errOut, err := fx.Run(cmd, names)
	expOut := expOutBuilder.String()

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
