package subaccount_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox/subaccount"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := subaccount.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "my-storage-box",
	}
	sbs := &hcloud.StorageBoxSubaccount{
		ID:         456,
		StorageBox: sb,
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSubaccountByID(gomock.Any(), sb, int64(456)).
		Return(sbs, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		DeleteSubaccount(gomock.Any(), sbs).
		Return(&hcloud.Action{ID: 789}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"my-storage-box", "456"})

	expOut := "Storage Box Subaccount 456 deleted\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := subaccount.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "my-storage-box",
	}

	snapshots := []*hcloud.StorageBoxSubaccount{
		{
			ID:         123,
			StorageBox: sb,
		},
		{
			ID:         456,
			StorageBox: sb,
		},
		{
			ID:         789,
			StorageBox: sb,
		},
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)

	var ids []string
	for _, sbs := range snapshots {
		ids = append(ids, strconv.FormatInt(sbs.ID, 10))
		fx.Client.StorageBoxClient.EXPECT().
			GetSubaccountByID(gomock.Any(), sb, sbs.ID).
			Return(sbs, nil, nil)
		fx.Client.StorageBoxClient.EXPECT().
			DeleteSubaccount(gomock.Any(), sbs).
			Return(&hcloud.Action{ID: sbs.ID + 1000}, nil, nil)
	}

	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(),
			&hcloud.Action{ID: 1123}, &hcloud.Action{ID: 1456}, &hcloud.Action{ID: 1789}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, append([]string{"my-storage-box"}, ids...))

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Storage Box Subaccounts 123, 456, 789 deleted\n", out)
}
