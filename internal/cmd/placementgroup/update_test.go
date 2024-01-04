package placementgroup

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateName(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.PlacementGroupClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.PlacementGroup{ID: 123}, nil, nil)
	fx.Client.PlacementGroupClient.EXPECT().
		Update(gomock.Any(), &hcloud.PlacementGroup{ID: 123}, hcloud.PlacementGroupUpdateOpts{
			Name: "new-name",
		})

	out, _, err := fx.Run(cmd, []string{"123", "--name", "new-name"})

	expOut := "placement group 123 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
