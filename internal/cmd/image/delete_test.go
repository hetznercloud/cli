package image

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

	image := &hcloud.Image{
		ID:   123,
		Name: "test",
	}

	fx.Client.ImageClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(image, nil, nil)
	fx.Client.ImageClient.EXPECT().
		Delete(gomock.Any(), image).
		Return(nil, nil)

	out, err := fx.Run(cmd, []string{"test"})

	expOut := "image test deleted\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
