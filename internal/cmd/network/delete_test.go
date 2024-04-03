package network_test

import (
	"fmt"
	"strings"
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

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := "Network test deleted\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := network.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	networks := []*hcloud.Network{
		{
			ID:   123,
			Name: "test1",
		},
		{
			ID:   456,
			Name: "test2",
		},
		{
			ID:   789,
			Name: "test3",
		},
	}

	expOutBuilder := strings.Builder{}

	var names []string
	for _, net := range networks {
		names = append(names, net.Name)
		expOutBuilder.WriteString(fmt.Sprintf("Network %s deleted\n", net.Name))
		fx.Client.NetworkClient.EXPECT().
			Get(gomock.Any(), net.Name).
			Return(net, nil, nil)
		fx.Client.NetworkClient.EXPECT().
			Delete(gomock.Any(), net).
			Return(nil, nil)
	}

	out, errOut, err := fx.Run(cmd, names)
	expOut := expOutBuilder.String()

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
