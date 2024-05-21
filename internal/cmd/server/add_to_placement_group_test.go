package server_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestAddToPlacementGroup(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.AddToPlacementGroupCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	var (
		server = hcloud.Server{
			ID:   42,
			Name: "my server",
		}
		placementGroup = hcloud.PlacementGroup{
			ID:   897,
			Name: "my Placement Group",
		}
	)

	fx.Client.Server.EXPECT().
		Get(gomock.Any(), server.Name).
		Return(&server, nil, nil)
	fx.Client.PlacementGroup.EXPECT().
		Get(gomock.Any(), placementGroup.Name).
		Return(&placementGroup, nil, nil)
	fx.Client.Server.EXPECT().
		AddToPlacementGroup(gomock.Any(), &server, &placementGroup)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), nil)

	out, errOut, err := fx.Run(cmd, []string{"-g", placementGroup.Name, server.Name})

	expOut := `Server 42 added to placement group my Placement Group
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
