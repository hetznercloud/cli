package image_test

import (
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		ID: 123,
	}

	fx.Client.ImageClient.EXPECT().
		GetByID(gomock.Any(), img.ID).
		Return(img, nil, nil)
	fx.Client.ImageClient.EXPECT().
		Delete(gomock.Any(), img).
		Return(nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"123"})

	expOut := "image 123 deleted\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := image.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	images := []*hcloud.Image{
		{
			ID: 123,
		},
		{
			ID: 456,
		},
		{
			ID: 789,
		},
	}

	var ids []string
	for _, img := range images {
		ids = append(ids, strconv.FormatInt(img.ID, 10))
		fx.Client.ImageClient.EXPECT().
			GetByID(gomock.Any(), img.ID).
			Return(img, nil, nil)
		fx.Client.ImageClient.EXPECT().
			Delete(gomock.Any(), img).
			Return(nil, nil)
	}

	out, errOut, err := fx.Run(cmd, ids)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "images 123, 456, 789 deleted\n", out)
}
