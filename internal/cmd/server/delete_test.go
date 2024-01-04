package server

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{
		ID:   123,
		Name: "test",
	}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(srv, nil, nil)
	fx.Client.ServerClient.EXPECT().
		DeleteWithResult(gomock.Any(), srv).
		Return(&hcloud.ServerDeleteResult{
			Action: &hcloud.Action{ID: 321},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321})

	out, _, err := fx.Run(cmd, []string{"test"})

	expOut := "Server test deleted\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
