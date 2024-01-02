package server_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestAddToPlacementGroup(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.AddToPlacementGroupCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
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

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), server.Name).
		Return(&server, nil, nil)
	fx.Client.PlacementGroupClient.EXPECT().
		Get(gomock.Any(), placementGroup.Name).
		Return(&placementGroup, nil, nil)
	fx.Client.ServerClient.EXPECT().
		AddToPlacementGroup(gomock.Any(), &server, &placementGroup)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), gomock.Any(), nil)

	out, _, err := fx.Run(cmd, []string{"-g", placementGroup.Name, server.Name})

	expOut := `Server 42 added to placement group my Placement Group
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
