package zone_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/zone"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := zone.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{
		ID:   123,
		Name: "example.com",
	}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		Delete(gomock.Any(), z).
		Return(hcloud.ZoneDeleteResult{Action: &hcloud.Action{ID: 42}}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 42}}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"123"})

	expOut := "Zone 123 deleted\n"

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := zone.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	zones := []*hcloud.Zone{
		{
			ID:   123,
			Name: "example1.com",
		},
		{
			ID:   456,
			Name: "example2.com",
		},
		{
			ID:   789,
			Name: "example3.com",
		},
	}

	var (
		names   []string
		actions []*hcloud.Action
	)
	for i, z := range zones {
		action := &hcloud.Action{ID: int64(i)}
		actions = append(actions, action)
		names = append(names, z.Name)
		fx.Client.ZoneClient.EXPECT().
			Get(gomock.Any(), z.Name).
			Return(z, nil, nil)
		fx.Client.ZoneClient.EXPECT().
			Delete(gomock.Any(), z).
			Return(hcloud.ZoneDeleteResult{Action: action}, nil, nil)
	}
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), actions).
		Return(nil)

	out, errOut, err := fx.Run(cmd, names)

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, "Zones example1.com, example2.com, example3.com deleted\n", out)
}
