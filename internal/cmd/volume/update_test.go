package volume

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

	fx.Client.VolumeClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Volume{ID: 123}, nil, nil)
	fx.Client.VolumeClient.EXPECT().
		Update(gomock.Any(), &hcloud.Volume{ID: 123}, hcloud.VolumeUpdateOpts{
			Name: "new-name",
		})

	out, err := fx.Run(cmd, []string{"123", "--name", "new-name"})

	expOut := "Volume 123 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
