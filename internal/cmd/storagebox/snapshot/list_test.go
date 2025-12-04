package snapshot_test

import (
	"fmt"
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

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := snapshot.ListCmd.CobraCommand(fx.State())

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "test",
	}
	sbs := &hcloud.StorageBoxSnapshot{
		ID:          456,
		Name:        "snapshot-1",
		Description: "some-description",
		Stats: hcloud.StorageBoxSnapshotStats{
			Size:           50 * util.Gibibyte,
			SizeFilesystem: 40 * util.Gibibyte,
		},
		IsAutomatic: false,
		Labels: map[string]string{
			"example.com/my": "label",
			"environment":    "prod",
			"just-a-key":     "",
		},
		Created:    time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
		StorageBox: sb,
	}

	fx.ExpectEnsureToken()
	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		AllSnapshotsWithOpts(
			gomock.Any(),
			sb,
			hcloud.StorageBoxSnapshotListOpts{Sort: []string{"id:asc"}},
		).
		Return([]*hcloud.StorageBoxSnapshot{
			sbs,
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := fmt.Sprintf(`ID    NAME         DESCRIPTION        SIZE     IS AUTOMATIC   AGE 
456   snapshot-1   some-description   50 GiB   no             %s
`, util.Age(sbs.Created, time.Now()))

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
