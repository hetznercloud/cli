package server_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestRemoveFromPlacementGroup(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.RemoveFromPlacementGroupCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := hcloud.Server{
		ID:   42,
		Name: "my server",
	}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), srv.Name).
		Return(&srv, nil, nil)
	fx.Client.ServerClient.EXPECT().
		RemoveFromPlacementGroup(gomock.Any(), &srv)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), nil)

	out, errOut, err := fx.Run(cmd, []string{srv.Name})

	expOut := `Server 42 removed from Placement Group
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
