package floatingip_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/floatingip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := floatingip.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	floatingIP := &hcloud.FloatingIP{
		ID:   123,
		Name: "test",
	}

	fx.Client.FloatingIPClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(floatingIP, nil, nil)
	fx.Client.FloatingIPClient.EXPECT().
		Delete(gomock.Any(), floatingIP).
		Return(nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := "Floating IP test deleted\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := floatingip.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	ips := []*hcloud.FloatingIP{
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

	var names []string
	for _, ip := range ips {
		names = append(names, ip.Name)
		fx.Client.FloatingIPClient.EXPECT().
			Get(gomock.Any(), ip.Name).
			Return(ip, nil, nil)
		fx.Client.FloatingIPClient.EXPECT().
			Delete(gomock.Any(), ip).
			Return(nil, nil)
	}

	out, errOut, err := fx.Run(cmd, names)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Floating IPs test1, test2, test3 deleted\n", out)
}
