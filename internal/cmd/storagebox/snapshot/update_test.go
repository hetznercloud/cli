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

func TestUpdate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := snapshot.UpdateCmd.CobraCommand(fx.State())
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
		UpdateSnapshot(gomock.Any(), sbs, hcloud.StorageBoxSnapshotUpdateOpts{
			Description: hcloud.Ptr("new description"),
		})

	out, errOut, err := fx.Run(cmd, []string{"my-storage-box", "my-snapshot", "--description", "new description"})

	expOut := "Storage Box Snapshot my-snapshot updated\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
