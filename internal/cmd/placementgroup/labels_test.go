package placementgroup_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/placementgroup"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.PlacementGroupClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.PlacementGroup{ID: 123}, nil, nil)
	fx.Client.PlacementGroupClient.EXPECT().
		Update(gomock.Any(), &hcloud.PlacementGroup{ID: 123}, hcloud.PlacementGroupUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label(s) key added to placement group 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := placementgroup.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	pg := &hcloud.PlacementGroup{
		ID: 123,
		Labels: map[string]string{
			"key": "value",
		},
	}

	fx.Client.PlacementGroupClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(pg, nil, nil)
	fx.Client.PlacementGroupClient.EXPECT().
		Update(gomock.Any(), pg, hcloud.PlacementGroupUpdateOpts{
			Labels: make(map[string]string),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label(s) key removed from placement group 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
