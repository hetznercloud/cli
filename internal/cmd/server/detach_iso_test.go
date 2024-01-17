package server_test

import (
	_ "embed"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDetachISO(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.DetachISOCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{ID: 123, Name: "my-server"}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)
	fx.Client.ServerClient.EXPECT().
		DetachISO(gomock.Any(), srv).
		Return(&hcloud.Action{ID: 789}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"my-server"}
	out, errOut, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "ISO detached from server 123\n", out)
}
