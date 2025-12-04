package snapshot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox/snapshot"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := snapshot.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "my-storage-box",
	}
	sbs := &hcloud.StorageBoxSnapshot{
		ID:         456,
		Name:       "my-snapshot",
		StorageBox: sb,
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSnapshot(gomock.Any(), sb, "my-snapshot").
		Return(sbs, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		DeleteSnapshot(gomock.Any(), sbs).
		Return(hcloud.StorageBoxSnapshotDeleteResult{Action: &hcloud.Action{ID: 789}}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"my-storage-box", "my-snapshot"})

	expOut := "Storage Box Snapshot my-snapshot deleted\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := snapshot.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "my-storage-box",
	}

	snapshots := []*hcloud.StorageBoxSnapshot{
		{
			ID:         123,
			Name:       "test1",
			StorageBox: sb,
		},
		{
			ID:         456,
			Name:       "test2",
			StorageBox: sb,
		},
		{
			ID:         789,
			Name:       "test3",
			StorageBox: sb,
		},
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)

	var names []string
	for _, sbs := range snapshots {
		names = append(names, sbs.Name)
		fx.Client.StorageBoxClient.EXPECT().
			GetSnapshot(gomock.Any(), sb, sbs.Name).
			Return(sbs, nil, nil)
		fx.Client.StorageBoxClient.EXPECT().
			DeleteSnapshot(gomock.Any(), sbs).
			Return(hcloud.StorageBoxSnapshotDeleteResult{Action: &hcloud.Action{ID: sbs.ID + 1000}}, nil, nil)
	}

	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(),
			&hcloud.Action{ID: 1123}, &hcloud.Action{ID: 1456}, &hcloud.Action{ID: 1789}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, append([]string{"my-storage-box"}, names...))

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Storage Box Snapshots test1, test2, test3 deleted\n", out)
}
