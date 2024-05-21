package server_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestReboot(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.RebootCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{ID: 123, Name: "my-server"}

	fx.Client.Server.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)
	fx.Client.Server.EXPECT().
		Reboot(gomock.Any(), srv).
		Return(&hcloud.Action{ID: 789}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"my-server"}
	out, errOut, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Server 123 rebooted\n", out)
}
