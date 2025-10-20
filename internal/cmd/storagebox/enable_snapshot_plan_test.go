package storagebox_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestEnableSnapshotPlan(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.EnableSnapshotPlanCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{ID: 123, Name: "my-storage-box"}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		EnableSnapshotPlan(gomock.Any(), sb, hcloud.StorageBoxEnableSnapshotPlanOpts{
			MaxSnapshots: 10,
			Minute:       0,
			Hour:         2,
			DayOfWeek:    hcloud.Ptr(time.Tuesday),
			DayOfMonth:   nil,
		}).
		Return(&hcloud.Action{ID: 456}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	args := []string{"my-storage-box", "--max-snapshots", "10", "--minute", "0", "--hour", "2", "--day-of-week", "tuesday"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Snapshot Plan enabled for Storage Box 123\n", out)
}
