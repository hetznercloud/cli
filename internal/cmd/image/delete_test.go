package image_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/image"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := image.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	img := &hcloud.Image{
		ID:   123,
		Name: "test",
	}

	fx.Client.ImageClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(img, nil, nil)
	fx.Client.ImageClient.EXPECT().
		Delete(gomock.Any(), img).
		Return(nil, nil)

	out, _, err := fx.Run(cmd, []string{"test"})

	expOut := "image test deleted\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
