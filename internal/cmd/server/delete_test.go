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

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.DeleteCmd.CobraCommand(fx.State())
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
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321})

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := "Server test deleted\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	servers := []*hcloud.Server{
		{
			ID:   123,
			Name: "test1",
		},
		{
			ID:   456,
			Name: "test2",
		},
		{
			ID:   789,
			Name: "test3",
		},
	}

	var (
		names   []string
		actions []*hcloud.Action
	)
	for i, srv := range servers {
		names = append(names, srv.Name)
		fx.Client.ServerClient.EXPECT().
			Get(gomock.Any(), srv.Name).
			Return(srv, nil, nil)
		fx.Client.ServerClient.EXPECT().
			DeleteWithResult(gomock.Any(), srv).
			Return(&hcloud.ServerDeleteResult{
				Action: &hcloud.Action{ID: int64(i)},
			}, nil, nil)
		actions = append(actions, &hcloud.Action{ID: int64(i)})
	}
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), actions)

	out, errOut, err := fx.Run(cmd, names)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Servers test1, test2, test3 deleted\n", out)
}
