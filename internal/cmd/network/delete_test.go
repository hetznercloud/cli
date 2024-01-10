package network_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/network"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := network.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	n := &hcloud.Network{
		ID:   123,
		Name: "test",
	}

	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(n, nil, nil)
	fx.Client.NetworkClient.EXPECT().
		Delete(gomock.Any(), n).
		Return(nil, nil)

	out, _, err := fx.Run(cmd, []string{"test"})

	expOut := "Network test deleted\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
