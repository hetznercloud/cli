package server_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

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
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321})

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := "Server test deleted\n"

	assert.NoError(t, err)
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

	expOutBuilder := strings.Builder{}

	var names []string
	for i, srv := range servers {
		names = append(names, srv.Name)
		expOutBuilder.WriteString(fmt.Sprintf("Server %s deleted\n", srv.Name))
		fx.Client.ServerClient.EXPECT().
			Get(gomock.Any(), srv.Name).
			Return(srv, nil, nil)
		fx.Client.ServerClient.EXPECT().
			DeleteWithResult(gomock.Any(), srv).
			Return(&hcloud.ServerDeleteResult{
				Action: &hcloud.Action{ID: int64(i)},
			}, nil, nil)
		fx.ActionWaiter.EXPECT().
			ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: int64(i)})
	}

	out, errOut, err := fx.Run(cmd, names)
	expOut := expOutBuilder.String()

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
