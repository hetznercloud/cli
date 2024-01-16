package primaryip_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/primaryip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := primaryip.DeleteCmd.CobraCommand(fx.State())
	ip := &hcloud.PrimaryIP{ID: 13}
	fx.ExpectEnsureToken()
	fx.Client.PrimaryIPClient.EXPECT().
		Get(
			gomock.Any(),
			"13",
		).
		Return(
			ip,
			&hcloud.Response{},
			nil,
		)
	fx.Client.PrimaryIPClient.EXPECT().
		Delete(
			gomock.Any(),
			ip,
		).
		Return(
			&hcloud.Response{},
			nil,
		)

	out, errOut, err := fx.Run(cmd, []string{"13"})

	expOut := "Primary IP 13 deleted\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
