package storagebox_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := storagebox.ListCmd.CobraCommand(fx.State())

	storageBox := &hcloud.StorageBox{
		ID:       123,
		Username: hcloud.Ptr("u12345"),
		Status:   hcloud.StorageBoxStatusActive,
		Name:     "test",
		Location: &hcloud.Location{Name: "fsn1"},
		Server:   hcloud.Ptr("u1337.your-storagebox.de"),
		System:   hcloud.Ptr("FSN1-BX355"),
		StorageBoxType: &hcloud.StorageBoxType{
			Name: "bx11",
		},
		Stats: &hcloud.StorageBoxStats{
			Size: 42 * util.Gibibyte,
		},
		Labels: map[string]string{
			"environment":    "prod",
			"example.com/my": "label",
			"just-a-key":     "",
		},
		Protection: hcloud.StorageBoxProtection{
			Delete: false,
		},
		SnapshotPlan: &hcloud.StorageBoxSnapshotPlan{
			MaxSnapshots: 10,
			Minute:       hcloud.Ptr(1),
			Hour:         hcloud.Ptr(2),
			DayOfWeek:    hcloud.Ptr(3),
			DayOfMonth:   hcloud.Ptr(4),
		},
		Created: time.Now().Add(-3 * time.Hour),
	}

	fx.ExpectEnsureToken()
	fx.Client.StorageBoxClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.StorageBoxListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
			},
		).
		Return([]*hcloud.StorageBox{
			storageBox,
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID    NAME   USERNAME   SERVER                     TYPE   SIZE     LOCATION   AGE
123   test   u12345     u1337.your-storagebox.de   bx11   42 GiB   fsn1       3h 
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
