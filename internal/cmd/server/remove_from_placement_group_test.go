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

func TestRemoveFromPlacementGroup(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.RemoveFromPlacementGroupCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	server := hcloud.Server{
		ID:   42,
		Name: "my server",
	}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), server.Name).
		Return(&server, nil, nil)
	fx.Client.ServerClient.EXPECT().
		RemoveFromPlacementGroup(gomock.Any(), &server)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), nil)

	out, _, err := fx.Run(cmd, []string{server.Name})

	expOut := `Server 42 removed from placement group
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
