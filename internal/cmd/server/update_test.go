package server_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateName(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Server{ID: 123}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Update(gomock.Any(), &hcloud.Server{ID: 123}, hcloud.ServerUpdateOpts{
			Name: "new-name",
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "--name", "new-name"})

	expOut := "Server 123 updated\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
