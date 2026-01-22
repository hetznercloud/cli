package snapshot_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox/snapshot"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := snapshot.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "my-storage-box",
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

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSnapshot(gomock.Any(), sb, "my-snapshot").
		Return(sbs, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"my-storage-box", "my-snapshot"})

	expOut := fmt.Sprintf(`ID:            456
Name:          snapshot-1
Description:   some-description
Created:       2024-01-02 15:04:05 UTC (%s)
Is automatic:  no

Stats:
  Size:             50 GiB
  Filesystem Size:  40 GiB

Labels:
  environment:     prod
  example.com/my:  label
  just-a-key:      

Storage Box:
  ID:  123
`, humanize.Time(sbs.Created))

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
