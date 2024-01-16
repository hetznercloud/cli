package volume_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/volume"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := volume.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	v := &hcloud.Volume{
		ID:   123,
		Name: "test",
	}

	fx.Client.VolumeClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(v, nil, nil)
	fx.Client.VolumeClient.EXPECT().
		Delete(gomock.Any(), v).
		Return(nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := "Volume test deleted\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
