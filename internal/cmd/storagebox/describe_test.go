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

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := storagebox.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	storageBox := &hcloud.StorageBox{
		ID:       123,
		Username: hcloud.Ptr("u12345"),
		Status:   hcloud.StorageBoxStatusActive,
		Name:     "test",
		StorageBoxType: &hcloud.StorageBoxType{
			ID:                     42,
			Name:                   "bx11",
			Description:            "BX11",
			SnapshotLimit:          hcloud.Ptr(10),
			AutomaticSnapshotLimit: hcloud.Ptr(10),
			SubaccountsLimit:       200,
			Size:                   1073741824,
			Pricings: []hcloud.StorageBoxTypeLocationPricing{
				{
					Location: "fsn1",
					PriceHourly: hcloud.Price{
						Net:   "1.0000",
						Gross: "1.1900",
					},
					PriceMonthly: hcloud.Price{
						Net:   "1.0000",
						Gross: "1.1900",
					},
					SetupFee: hcloud.Price{
						Net:   "1.0000",
						Gross: "1.1900",
					},
				},
			},
			DeprecatableResource: hcloud.DeprecatableResource{
				Deprecation: &hcloud.DeprecationInfo{
					UnavailableAfter: time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
					Announced:        time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		Location: &hcloud.Location{
			ID:          42,
			Name:        "fsn1",
			Description: "Falkenstein DC Park 1",
			Country:     "DE",
			City:        "Falkenstein",
			Latitude:    50.47612,
			Longitude:   12.370071,
			NetworkZone: "eu-central",
		},
		AccessSettings: hcloud.StorageBoxAccessSettings{
			ReachableExternally: false,
			SambaEnabled:        false,
			SSHEnabled:          false,
			WebDAVEnabled:       false,
			ZFSEnabled:          false,
		},
		Server: hcloud.Ptr("u1337.your-storagebox.de"),
		System: hcloud.Ptr("FSN1-BX355"),
		Stats: &hcloud.StorageBoxStats{
			Size:          0,
			SizeData:      0,
			SizeSnapshots: 0,
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
			MaxSnapshots: 0,
			Minute:       nil,
			Hour:         nil,
			DayOfWeek:    nil,
			DayOfMonth:   nil,
		},
		Created: time.Date(2016, 1, 30, 23, 55, 0, 0, time.UTC),
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(storageBox, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:							123
Name:						test
Created:					Sat Jan 30 23:55:00 UTC 2016 (9 years ago)
Status:						active
Username:					u12345
Server:						u1337.your-storagebox.de
System:						FSN1-BX355
Snapshot Plan:
  Max Snapshots:			0
Protection:
  Delete:					false
Stats:
  Size:						0 B
  Size Data:				0 B
  Size Snapshots:			0 B
Labels:
  environment: prod
  example.com/my: label
  just-a-key: 
Access Settings:
  Reachable Externally:		false
  Samba Enabled:			false
  SSH Enabled:				false
  WebDAV Enabled:			false
  ZFS Enabled:				false
Storage Box Type:
  ID:						42
  Name:						bx11
  Description:				BX11
  Size:						1.0 GiB
  Snapshot Limit:			10
  Automatic Snapshot Limit:	10
  Subaccounts Limit:		200
Location:
  ID:						42
  Name:						fsn1
  Description:				Falkenstein DC Park 1
  Network Zone:				eu-central
  Country:					DE
  City:						Falkenstein
  Latitude:					50.476120
  Longitude:				12.370071
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
