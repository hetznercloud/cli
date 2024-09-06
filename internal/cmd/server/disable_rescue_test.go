package server_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDisableRescue(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.DisableRescueCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{ID: 123, Name: "my-server"}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)
	fx.Client.ServerClient.EXPECT().
		DisableRescue(gomock.Any(), srv).
		Return(&hcloud.Action{ID: 789}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"my-server"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Rescue disabled for server 123\n", out)
}
