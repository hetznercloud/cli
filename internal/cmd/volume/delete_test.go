package volume

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DeleteCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	volume := &hcloud.Volume{
		ID:   123,
		Name: "test",
	}

	fx.Client.VolumeClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(volume, nil, nil)
	fx.Client.VolumeClient.EXPECT().
		Delete(gomock.Any(), volume).
		Return(nil, nil)

	out, err := fx.Run(cmd, []string{"test"})

	expOut := "Volume test deleted\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
