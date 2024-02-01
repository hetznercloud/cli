package server_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestChangeType(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.ChangeTypeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{ID: 123, Name: "my-server"}
	st := &hcloud.ServerType{ID: 456, Name: "cax21"}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)
	fx.Client.ServerTypeClient.EXPECT().
		Get(gomock.Any(), "cax21").
		Return(st, nil, nil)
	fx.Client.ServerClient.EXPECT().
		ChangeType(gomock.Any(), srv, hcloud.ServerChangeTypeOpts{
			ServerType:  st,
			UpgradeDisk: true,
		}).
		Return(&hcloud.Action{ID: 789}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"my-server", "cax21"}
	out, errOut, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Server 123 changed to type cax21\n", out)
}

func TestChangeTypeKeepDisk(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.ChangeTypeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{ID: 123, Name: "my-server"}
	st := &hcloud.ServerType{ID: 456, Name: "cax21"}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)
	fx.Client.ServerTypeClient.EXPECT().
		Get(gomock.Any(), "cax21").
		Return(st, nil, nil)
	fx.Client.ServerClient.EXPECT().
		ChangeType(gomock.Any(), srv, hcloud.ServerChangeTypeOpts{
			ServerType:  st,
			UpgradeDisk: false,
		}).
		Return(&hcloud.Action{ID: 789}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"my-server", "cax21", "--keep-disk"}
	out, errOut, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Server 123 changed to type cax21 (disk size was unchanged)\n", out)
}
