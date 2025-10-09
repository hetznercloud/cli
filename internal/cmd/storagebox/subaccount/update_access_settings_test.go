package subaccount_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox/subaccount"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateAccessSettings(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := subaccount.UpdateAccessSettingsCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{ID: 123, Name: "my-storage-box"}
	sbs := &hcloud.StorageBoxSubaccount{ID: 456, StorageBox: sb}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSubaccount(gomock.Any(), sb, "456").
		Return(sbs, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		UpdateSubaccountAccessSettings(gomock.Any(), sbs, hcloud.StorageBoxSubaccountAccessSettingsUpdateOpts{
			SambaEnabled:        nil,
			SSHEnabled:          hcloud.Ptr(true),
			WebDAVEnabled:       hcloud.Ptr(false),
			ReachableExternally: nil,
			Readonly:            hcloud.Ptr(true),
		}).
		Return(&hcloud.Action{ID: 456}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	args := []string{"my-storage-box", "456", "--enable-ssh", "--enable-webdav=false", "--readonly=true"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Access settings updated for Storage Box Subaccount 456\n", out)
}
