package storagebox_test

import (
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/create_response.json
var createResponseJSON string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.StorageBoxClient.EXPECT().
		Create(gomock.Any(), hcloud.StorageBoxCreateOpts{
			Name:           "my-storage-box",
			StorageBoxType: &hcloud.StorageBoxType{Name: "bx11"},
			Location:       &hcloud.Location{Name: "fsn1"},
			Password:       "my-password",
			AccessSettings: &hcloud.StorageBoxCreateOptsAccessSettings{
				SambaEnabled:        hcloud.Ptr(true),
				SSHEnabled:          hcloud.Ptr(true),
				ZFSEnabled:          hcloud.Ptr(true),
				ReachableExternally: hcloud.Ptr(false),
				WebDAVEnabled:       hcloud.Ptr(false),
			},
			Labels:  make(map[string]string),
			SSHKeys: []string{"mykey"},
		}).
		Return(hcloud.StorageBoxCreateResult{
			StorageBox: &hcloud.StorageBox{
				ID:   123,
				Name: "my-storage-box",
			},
			Action: &hcloud.Action{ID: 456},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"--name", "my-storage-box", "--type", "bx11", "--location", "fsn1",
		"--password", "my-password", "--enable-samba", "--enable-ssh", "--enable-zfs", "--ssh-key", "mykey"})

	expOut := "Storage Box 123 created\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.StorageBoxClient.EXPECT().
		Create(gomock.Any(), hcloud.StorageBoxCreateOpts{
			Name:           "my-storage-box",
			StorageBoxType: &hcloud.StorageBoxType{Name: "bx11"},
			Location:       &hcloud.Location{Name: "fsn1"},
			Password:       "my-password",
			AccessSettings: &hcloud.StorageBoxCreateOptsAccessSettings{
				SambaEnabled:        hcloud.Ptr(true),
				SSHEnabled:          hcloud.Ptr(true),
				ZFSEnabled:          hcloud.Ptr(true),
				ReachableExternally: hcloud.Ptr(false),
				WebDAVEnabled:       hcloud.Ptr(false),
			},
			Labels:  make(map[string]string),
			SSHKeys: []string{"mykey"},
		}).
		Return(hcloud.StorageBoxCreateResult{
			StorageBox: &hcloud.StorageBox{
				ID:       42,
				Username: hcloud.Ptr("u12345"),
				Status:   hcloud.StorageBoxStatusActive,
				Name:     "string",
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
							Announced:        time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
							UnavailableAfter: time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
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
			},
			Action: &hcloud.Action{ID: 456},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name", "my-storage-box", "--type", "bx11", "--location", "fsn1",
		"--password", "my-password", "--enable-samba", "--enable-ssh", "--enable-zfs", "--ssh-key", "mykey"})

	expOut := "Storage Box 42 created\n"

	require.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJSON, jsonOut)
}

func TestCreateProtection(t *testing.T) {
	// TODO implement once change-protection is implemented
}
