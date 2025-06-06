package primaryip_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

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

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := primaryip.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	ips := []*hcloud.PrimaryIP{
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
		fx.Client.PrimaryIPClient.EXPECT().
			Get(gomock.Any(), ip.Name).
			Return(ip, nil, nil)
		fx.Client.PrimaryIPClient.EXPECT().
			Delete(gomock.Any(), ip).
			Return(nil, nil)
	}

	out, errOut, err := fx.Run(cmd, names)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Primary IPs test1, test2, test3 deleted\n", out)
}
