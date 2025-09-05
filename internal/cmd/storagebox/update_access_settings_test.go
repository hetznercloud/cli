package storagebox_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateAccessSettings(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.UpdateAccessSettingsCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{ID: 123, Name: "my-storage-box"}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		UpdateAccessSettings(gomock.Any(), sb, hcloud.StorageBoxUpdateAccessSettingsOpts{
			SambaEnabled:        nil,
			SSHEnabled:          hcloud.Ptr(true),
			WebDAVEnabled:       hcloud.Ptr(false),
			ZFSEnabled:          hcloud.Ptr(true),
			ReachableExternally: nil,
		}).
		Return(&hcloud.Action{ID: 456}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	args := []string{"my-storage-box", "--ssh-enabled", "--webdav-enabled=false", "--zfs-enabled=true"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Access settings updated for Storage Box 123\n", out)
}
