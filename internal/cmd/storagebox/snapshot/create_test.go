package snapshot_test

import (
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox/snapshot"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/create_response.json
var createResponseJSON string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := snapshot.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "my-storage-box",
	}
	sbs := &hcloud.StorageBoxSnapshot{
		ID:   456,
		Name: "snapshot-1",
		Stats: &hcloud.StorageBoxSnapshotStats{
			Size: 50 * util.Gibibyte,
		},
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		CreateSnapshot(gomock.Any(), sb, hcloud.StorageBoxSnapshotCreateOpts{
			Description: "some-description",
			Labels:      map[string]string{"foo": "bar"},
		}).
		Return(hcloud.StorageBoxSnapshotCreateResult{
			Snapshot: sbs,
			Action:   &hcloud.Action{ID: 789},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSnapshotByID(gomock.Any(), sb, sbs.ID).
		Return(sbs, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"--description", "some-description", "--label", "foo=bar", "my-storage-box"})

	expOut := `Storage Box Snapshot 456 created
Name: snapshot-1
Size: 50 GiB
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := snapshot.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "my-storage-box",
	}
	sbs := &hcloud.StorageBoxSnapshot{
		ID:          456,
		Name:        "snapshot-1",
		Description: "some-description",
		Stats: &hcloud.StorageBoxSnapshotStats{
			Size:           50 * util.Gibibyte,
			SizeFilesystem: 40 * util.Gibibyte,
		},
		IsAutomatic: false,
		Labels: map[string]string{
			"example.com/my": "label",
			"environment":    "prod",
			"just-a-key":     "",
		},
		Created: time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
		StorageBox: &hcloud.StorageBox{
			ID: 1337,
		},
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		CreateSnapshot(gomock.Any(), sb, hcloud.StorageBoxSnapshotCreateOpts{
			Description: "some-description",
			Labels:      map[string]string{},
		}).
		Return(hcloud.StorageBoxSnapshotCreateResult{
			Snapshot: sbs,
			Action:   &hcloud.Action{ID: 789},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSnapshotByID(gomock.Any(), sb, sbs.ID).
		Return(sbs, nil, nil)

	jsonOut, out, err := fx.Run(cmd, []string{"--description", "some-description", "my-storage-box", "-o=json"})

	expOut := "Storage Box Snapshot 456 created\n"

	require.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJSON, jsonOut)
}
