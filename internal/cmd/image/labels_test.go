package image_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/image"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := image.LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ImageClient.EXPECT().
		GetByID(gomock.Any(), int64(123)).
		Return(&hcloud.Image{ID: 123}, nil, nil)
	fx.Client.ImageClient.EXPECT().
		Update(gomock.Any(), &hcloud.Image{ID: 123}, hcloud.ImageUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label(s) key added to image 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := image.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ImageClient.EXPECT().
		GetByID(gomock.Any(), int64(123)).
		Return(&hcloud.Image{
			ID: 123,
			Labels: map[string]string{
				"key": "value",
			},
		}, nil, nil)
	fx.Client.ImageClient.EXPECT().
		Update(gomock.Any(), &hcloud.Image{ID: 123}, hcloud.ImageUpdateOpts{
			Labels: make(map[string]string),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label(s) key removed from image 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
