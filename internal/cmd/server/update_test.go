package server

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateName(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := UpdateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Server{ID: 123}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Update(gomock.Any(), &hcloud.Server{ID: 123}, hcloud.ServerUpdateOpts{
			Name: "new-name",
		})

	out, _, err := fx.Run(cmd, []string{"123", "--name", "new-name"})

	expOut := "Server 123 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
