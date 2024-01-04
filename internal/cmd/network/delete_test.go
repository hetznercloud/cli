package network

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	network := &hcloud.Network{
		ID:   123,
		Name: "test",
	}

	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(network, nil, nil)
	fx.Client.NetworkClient.EXPECT().
		Delete(gomock.Any(), network).
		Return(nil, nil)

	out, _, err := fx.Run(cmd, []string{"test"})

	expOut := "Network test deleted\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
