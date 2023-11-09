package placementgroup

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := LabelCmds.AddCobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
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

	out, _, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label key added to placement group 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := LabelCmds.RemoveCobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.PlacementGroupClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.PlacementGroup{
			ID: 123,
			Labels: map[string]string{
				"key": "value",
			},
		}, nil, nil)
	fx.Client.PlacementGroupClient.EXPECT().
		Update(gomock.Any(), &hcloud.PlacementGroup{ID: 123}, hcloud.PlacementGroupUpdateOpts{
			Labels: make(map[string]string),
		})

	out, _, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label key removed from placement group 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
