package storagebox_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestRollbackSnapshot(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.RollbackSnapshotCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{ID: 123, Name: "my-storage-box"}
	sbs := &hcloud.StorageBoxSnapshot{StorageBox: sb, ID: 456, Name: "my-snapshot"}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSnapshot(gomock.Any(), sb, "my-snapshot").
		Return(sbs, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		RollbackSnapshot(gomock.Any(), sb, hcloud.StorageBoxRollbackSnapshotOpts{Snapshot: sbs}).
		Return(&hcloud.Action{ID: 789}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"my-storage-box", "--snapshot", "my-snapshot"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Rolled back Storage Box 123 to Snapshot 456\n", out)
}
