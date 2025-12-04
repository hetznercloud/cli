package subaccount_test

import (
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox/subaccount"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/create_response.json
var createResponseJSON string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := subaccount.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "my-storage-box",
	}
	sbs := &hcloud.StorageBoxSubaccount{
		ID:         1,
		StorageBox: sb,
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		CreateSubaccount(gomock.Any(), sb, hcloud.StorageBoxSubaccountCreateOpts{
			HomeDirectory: "/home/directory",
			Password:      "my-password",
			AccessSettings: &hcloud.StorageBoxSubaccountCreateOptsAccessSettings{
				Readonly:   hcloud.Ptr(true),
				SSHEnabled: hcloud.Ptr(true),
			},
			Labels: make(map[string]string),
		}).
		Return(hcloud.StorageBoxSubaccountCreateResult{
			Subaccount: sbs,
			Action:     &hcloud.Action{ID: 456},
		}, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSubaccountByID(gomock.Any(), sb, sbs.ID).
		Return(sbs, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"--password", "my-password", "--home-directory", "/home/directory",
		"--enable-ssh", "--readonly=true", "my-storage-box"})

	expOut := "Storage Box Subaccount 1 created\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := subaccount.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   42,
		Name: "my-storage-box",
	}

	sbs := &hcloud.StorageBoxSubaccount{
		ID:            42,
		Name:          "u1337-sub1",
		Username:      "u1337-sub1",
		HomeDirectory: "my_backups/host01.my.company",
		Server:        "u1337-sub1.your-storagebox.de",
		AccessSettings: &hcloud.StorageBoxSubaccountAccessSettings{
			ReachableExternally: false,
			Readonly:            false,
			SambaEnabled:        false,
			SSHEnabled:          false,
			WebDAVEnabled:       false,
		},
		Description: "host01 backup",
		Labels: map[string]string{
			"environment":    "prod",
			"example.com/my": "label",
			"just-a-key":     "",
		},
		Created:    time.Date(2016, 1, 30, 23, 55, 0, 0, time.UTC),
		StorageBox: sb,
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		CreateSubaccount(gomock.Any(), sb, hcloud.StorageBoxSubaccountCreateOpts{
			HomeDirectory: "/home/directory",
			Password:      "my-password",
			AccessSettings: &hcloud.StorageBoxSubaccountCreateOptsAccessSettings{
				Readonly:   hcloud.Ptr(true),
				SSHEnabled: hcloud.Ptr(true),
			},
			Labels: make(map[string]string),
		}).
		Return(hcloud.StorageBoxSubaccountCreateResult{
			Subaccount: sbs,
			Action:     &hcloud.Action{ID: 456},
		}, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSubaccountByID(gomock.Any(), sb, sbs.ID).
		Return(sbs, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--password", "my-password", "--home-directory", "/home/directory",
		"--enable-ssh", "--readonly=true", "my-storage-box"})

	expOut := "Storage Box Subaccount 42 created\n"

	require.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJSON, jsonOut)
}
