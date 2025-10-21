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

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := snapshot.LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "test",
	}
	sbs := &hcloud.StorageBoxSnapshot{
		ID:         456,
		Name:       "my-snapshot",
		StorageBox: sb,
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSnapshot(gomock.Any(), sb, "my-snapshot").
		Return(sbs, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		UpdateSnapshot(gomock.Any(), sbs, hcloud.StorageBoxSnapshotUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "my-snapshot", "key=value"})

	expOut := "Label(s) key added to Storage Box Snapshot my-snapshot\n"

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := snapshot.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "test",
	}
	sbs := &hcloud.StorageBoxSnapshot{
		ID:   456,
		Name: "my-snapshot",
		Labels: map[string]string{
			"key": "value",
		},
		StorageBox: sb,
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSnapshot(gomock.Any(), sb, "my-snapshot").
		Return(sbs, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		UpdateSnapshot(gomock.Any(), sbs, hcloud.StorageBoxSnapshotUpdateOpts{
			Labels: make(map[string]string),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "my-snapshot", "key"})

	expOut := "Label(s) key removed from Storage Box Snapshot my-snapshot\n"

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}
